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
    pname = "nopresh-server";
    version = "0.1.0";
    src = ./server;
    vendorHash = "sha256-Rmdw2k5oseSfa4SUp+Lsyc2Vv9rC20dSa2VvzuPlHiw=";
    subPackages = ["cmd/api"];
    env.CGO_ENABLED = "0";
    # Embed the tz database into the binary so TimeZone=... in the DSN resolves
    # in a minimal image that has no /usr/share/zoneinfo.
    tags = ["timetzdata"];
    doCheck = false;
  };
in {
  # https://devenv.sh/basics/

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    nixd

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

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  # Build + import the server image straight into rootless podman storage in one
  # step (skips the docker-archive tarball + `podman load`). Run: `load-server`.
  scripts.load-server.exec = ''
    devenv container copy --registry \
      "containers-storage:[overlay@$HOME/.local/share/containers/storage+''${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/containers]" \
      server
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

  # processes.build-backend = {
  #   exec = ''
  #     go build
  #   '';
  # };

  containers = {
    server = {
      name = "nopresh-server";
      enableLayerDeduplication = true;
      # Exec the static binary directly. The default entrypoint sources the whole
      # dev shell (`shell.envScript`), which drags the entire toolchain closure
      # (~7GB) into the image; bypassing it keeps only the binary + base utils.
      entrypoint = ["${backend}/bin/api"];
      startupCommand = [];
      # Only ship the backend package, not the whole repo (the default).
      copyToRoot = [backend];
    };
  };

  # env = lib.mkMerge [
  
  # ];

  env = lib.mkMerge [
    (let
      apiPort = toString config.processes.backend.ports.http.value;
      clientPort = toString config.processes.client.ports.http.value;
      database = builtins.head config.services.postgres.initialDatabases;
      databasePort = toString config.services.postgres.port;
    in {
      # Empty host => listen on all interfaces (dual-stack IPv4 + IPv6) so that
      # a client using "localhost" reaches the server whether localhost resolves
      # to 127.0.0.1 or ::1.
      API_HOST = "0.0.0.0";
      API_PORT = apiPort;
      DOMAINS = "http://localhost:${clientPort}";
      ENVIRONMENT = "development";
      JWT_SECRET_KEY = "supersecretkey";
      JWT_DURATION = "15m";
      DB_HOST = "localhost";
      DB_PORT = databasePort;
      DB_USER = database.user;
      DB_PASSWORD = database.pass;
      DB_NAME = database.name;
      # Not `TZ`: that's a reserved libc var and gets clobbered by the shell;
      # use an app-specific name the server reads explicitly.
      APP_TZ = "Europe/Zagreb";
      VITE_BASE_URL = "http://localhost:${apiPort}/api";
      VITE_PORT = lib.toInt clientPort;
    })
    # devenv bakes the whole `env` into container images. During a container
    # build (`container.isBuilding`), blank vars we don't want in the image:
    #  - big store paths (dev profile + Go toolchain) -> keeps the image small
    #  - dev-only secrets/config -> must be supplied at runtime (`podman/docker/whatever run -e`)
    # The dev shell is unaffected; the runtime binary needs none of these baked.
    (lib.mkIf config.container.isBuilding {
      DEVENV_PROFILE = lib.mkForce "";
      DEVENV_TASK_FILE = lib.mkForce "";
      GOROOT = lib.mkForce "";
      JWT_SECRET_KEY = lib.mkForce "";
      NOPRESH_DB_DSN = lib.mkForce "";
    })
  ];

  # See full reference at https://devenv.sh/reference/options/
}
