{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git

    # Protobuf stuff
    protols
    buf
    pb
    protoc-gen-connect-go
    protoc-gen-es

    bruno

    # Go
    gorm-gentool

    vitejs
    tailwindcss
    vscode-langservers-extracted
    superhtml
  ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
    version = "1.26.4";
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
  # processes.dev.exec = "${lib.getExe pkgs.watchexec} -n -- ls -la";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  # https://devenv.sh/basics/
  enterShell = ''
    hello         # Run scripts directly
    git --version # Use packages
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  env.NOPRESH_DB_DSN = "host=localhost user=user password=password dbname=nopresh port=5433 sslmode=disable TimeZone=Europe/Zagreb";

  
  # See full reference at https://devenv.sh/reference/options/
}
