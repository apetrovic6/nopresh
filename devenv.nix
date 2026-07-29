{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  # Build the backend as a standalone, static Nix package (no CGO -> runs in a
  # minimal image). The binary is named after its subpackage dir: cmd/api -> "api".
  backend = (pkgs.buildGoModule.override {go = config.languages.go.package;}) {
    pname = "nopresh";
    version = "0.1.0";
    src = lib.cleanSourceWith {
      src = ./server;
      filter = path: _type: !lib.hasInfix "cmd/api/dist/" (toString path + "/");
    };
    vendorHash = "sha256-Rmdw2k5oseSfa4SUp+Lsyc2Vv9rC20dSa2VvzuPlHiw=";
    subPackages = ["cmd/api"];
    env.CGO_ENABLED = "0";
    # Embed the tz database into the binary so TimeZone=... in the DSN resolves
    # in a minimal image that has no /usr/share/zoneinfo.
    tags = ["timetzdata"];
    doCheck = false;
    preBuild = ''
      mkdir -p cmd/api/dist
      cp -r ${frontend}/. cmd/api/dist/
    '';
  };

  # Build the frontend (TanStack Start SSR) reproducibly with pnpm. `pnpm build`
  # emits a self-contained `.output` (nitro `bun` preset) that we run with bun.
  frontend = pkgs.stdenv.mkDerivation (finalAttrs: {
    pname = "nopresh-web";
    version = "0.1.0";
    # Exclude local node_modules / build outputs so they aren't copied into the
    # store (bloat + pnpm's non-interactive purge error).
    src = lib.cleanSourceWith {
      src = ./web;
      filter = path: _type: let
        b = baseNameOf path;
      in
        b != "node_modules" && b != ".output" && b != "dist";
    };
    nativeBuildInputs = [
      pkgs.nodejs
      pkgs.pnpm.configHook
    ];
    pnpmDeps = pkgs.pnpm.fetchDeps {
      inherit (finalAttrs) pname version src;
      # pnpm 11 requires the newer fetcher (v3 is unsupported; 4 is latest).
      fetcherVersion = 4;
      hash = "sha256-HexFVoiEp+e1HyEZXRZVqOePqxwlatK7ODYPWLOgTvQ=";
    };
    # CI=true so pnpm runs non-interactively (no TTY in the sandbox).
    # VITE_-prefixed vars are inlined into the browser bundle at build time.
    env.CI = "true";
    buildPhase = ''
      runHook preBuild
      pnpm build
      runHook postBuild
    '';
    installPhase = ''
      runHook preInstall
      # SPA build output static files under dist/client
      cp -r dist/client $out
      runHook postInstall
    '';
  });
in {
  # https://devenv.sh/basics/

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    nixd

    podman-compose

    # Protobuf stuff
    protols
    buf
    pb
    protoc-gen-connect-go
    protoc-gen-es

    bruno

    # Go
    gorm-gentool
    vscode-langservers-extracted
    superhtml
  ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
    version = "1.26.5";
  };

  languages.typescript = {
    enable = true;
  };

  languages.javascript = {
    enable = true;
    bun.enable = true;
    pnpm.enable = true;
  };

  services.postgres = {
    enable = true;
    listen_addresses = "0.0.0.0";
    port = 5433;
    createDatabase = true;
    package = pkgs.postgresql_18;
    initialDatabases = let
      user = "user";
      pass = "password";
    in [
      {
        name = "nopresh";
        inherit user pass;
      }
    ];
  };

  services.adminer = {
    enable = true;
    listen = "127.0.0.1:8080";
  };

  # https://devenv.sh/processes/
  # Regenerate the go + typescript clients whenever a .proto file changes.
  processes.protoWatch.exec = ''
    ${lib.getExe pkgs.watchexec} \
      --watch ./server/proto \
      --exts proto \
      --on-busy-update restart \
      -- devenv tasks run nopresh:rebuildServerProto nopresh:rebuildClientProto
  '';

  # import the server image straight into rootless podman storage in one
  scripts.load-server.exec = ''
    devenv container copy --registry \
      "containers-storage:[overlay@$HOME/.local/share/containers/storage+''${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/containers]" \
      nopresh
  '';

  # https://devenv.sh/basics/
  enterShell = ''
    hello         # Run scripts directly
    git --version # Use packages
  '';

  # https://devenv.sh/tasks/
  tasks = {
    "nopresh:rebuildServerProto" = {
      exec = "cd ./server && buf lint && buf generate";
      execIfModified = ["./server/proto/**/*.proto"];
      before = ["devenv:enterShell"];
    };

    "nopresh:rebuildClientProto" = {
      exec = "cd ./web && pnpm generate";
      execIfModified = ["./server/proto/**/*.proto"];
      before = ["devenv:enterShell"];
    };
  };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  processes.backend = {
    cwd = "${config.git.root}/server";
    exec = ''
      echo "HTTP server on port ${toString config.processes.backend.ports.http.value}"
      go run ./cmd/api/
    '';
    after = ["devenv:processes:postgres"];
    # NB: avoid browser-blocked "unsafe" ports (6000=X11, 6666/6667=IRC, etc.);
    # browsers reject fetches to them with net::ERR_UNSAFE_PORT.
    ports.http.allocate = 5000;
  };

  processes.client = {
    cwd = "${config.git.root}/web";
    exec = ''
      bun dev
    '';
    after = ["devenv:processes:backend"];
    ports.http.allocate = 3000;
  };

  containers.nopresh = {
    name = "nopresh";
    enableLayerDeduplication = true;
    # Exec the static binary directly. The default entrypoint sources the whole
    # dev shell (`shell.envScript`), which drags the entire toolchain closure
    # (~7GB) into the image; bypassing it keeps only the binary + base utils.
    entrypoint = ["${backend}/bin/api"];
    startupCommand = [];
    # Only ship the backend package, not the whole repo (the default).
    copyToRoot = [backend];
  };

  env = lib.mkMerge [
    (let
      apiHost = "0.0.0.0";
      apiPort = toString config.processes.backend.ports.http.value;
    in {
      # Safe to bake into the image (no secrets, no big store paths).
      # Empty host => listen on all interfaces (dual-stack IPv4 + IPv6) so that
      # a client using "localhost" reaches the server whether localhost resolves
      # to 127.0.0.1 or ::1.
      HOST = apiHost;
      PORT = apiPort;
      JWT_DURATION = "15m";
      # Not `TZ`: that's a reserved libc var and gets clobbered by the shell;
      # use an app-specific name the server reads explicitly.
      APP_TZ = "Europe/Zagreb";
      ENVIRONMENT =
        if config.container.isBuilding
        then "production"
        else "development";
    })

    # Dev-only secrets + local service wiring. Defined ONLY when NOT building a
    # container, so they can never be baked into the image (devenv bakes the whole
    # `env`). Supply them at runtime instead (`podman/docker run -e ...`). This is
    # fail-safe: a new secret added here is excluded from images by default.
    (lib.mkIf (!config.container.isBuilding) (let
      database = builtins.head config.services.postgres.initialDatabases;
      databasePort = toString config.services.postgres.port;
    in {
      JWT_SECRET_KEY = "supersecretkey";
      DB_HOST = "localhost";
      DB_PORT = databasePort;
      DB_USER = database.user;
      DB_PASSWORD = database.pass;
      DB_NAME = database.name;
    }))

    # These are set by other modules (dev profile, task runner, Go toolchain) and
    # point at huge store paths, so they can't just be "left undefined" — blank
    # them during container builds to keep the image small.
    (lib.mkIf config.container.isBuilding {
      DEVENV_PROFILE = lib.mkForce "";
      DEVENV_TASK_FILE = lib.mkForce "";
      GOROOT = lib.mkForce "";
    })
  ];

  # See full reference at https://devenv.sh/reference/options/
}
