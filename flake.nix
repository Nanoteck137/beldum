{
  description = "Devshell for beldum";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";

    tagopus.url      = "github:nanoteck137/tagopus/v0.1.1";
    tagopus.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, devtools, tagopus, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        beldum = pkgs.buildGoModule {
          pname = "beldum";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/beldum" "cmd/beldum-cli" "cmd/beldum-migrate"];

          ldflags = [
            "-X github.com/nanoteck137/beldum.Version=${version}"
            "-X github.com/nanoteck137/beldum.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          tags = ["fts5"];

          vendorHash = "sha256-tW/hbgKfTEDzknssIozaRNpo8vVpdSQF1jiHT2QyA3w=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/beldum --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
            wrapProgram $out/bin/beldum-migrate --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
            wrapProgram $out/bin/beldum-cli --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick tagopus.packages.${system}.default ]}
          '';
        };

        beldumWeb = pkgs.buildNpmPackage {
          name = "beldum-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-aR/7JKL10c9QgD1hpke5kV0wnHBYXaGp4M0hJwi0CAI=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            echo '{ "type": "module" }' > $out/package.json

            mkdir $out/bin
            echo -e "#!${pkgs.runtimeShell}\n${pkgs.nodejs}/bin/node $out\n" > $out/bin/beldum-web
            chmod +x $out/bin/beldum-web

            runHook postInstall
          '';
        };

        tools = devtools.packages.${system};
      in
      {
        packages = {
          default = beldum;
          beldum = beldum;
          beldum-web = beldumWeb;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            imagemagick
            ffmpeg

            tagopus.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.default = import ./nix/beldum.nix { inherit self; };
      nixosModules.beldum-web = import ./nix/beldum-web.nix { inherit self; };
    };
}
