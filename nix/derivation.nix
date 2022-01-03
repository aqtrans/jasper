{ stdenv, lib, buildGoModule, git, makeWrapper, substituteAll, grpc, protobuf, pkgs, runCommandLocal }:

buildGoModule rec {
  pname = "jasper";
  #version = "0.0.1";

  # dynamic version based on git; https://blog.replit.com/nix_dynamic_version
  revision = runCommandLocal "get-rev" {
          nativeBuildInputs = [ git ];
      } "GIT_DIR=${src}/.git git rev-parse --short HEAD | tr -d '\n' > $out";  

  buildDate = runCommandLocal "get-date" {} "date +'%Y-%m-%d_%T' | tr -d '\n' > $out"; 

  version = "0" + builtins.readFile revision;        

  src = ../.;

  nativeBuildInputs = [ 
    makeWrapper 
  ];

  buildInputs = [ 
    git 
  ];

  ldflags = [ "-X main.sha1ver=${builtins.readFile revision}" "-X main.buildTime=${builtins.readFile buildDate}" ];

  vendorSha256 = "02cfxd4508j481jdpbaa9s2nc2n9hzcswlwsq1fiyfvl492v4djg";

  runVend = false;

  deleteVendor = false;

  subPackages = [ "./" ];

  meta = with lib; {
    description = "A meme as a service? That's a Paddlin'";
    homepage = "https://github.com/aqtrans/thatsapaddlin";
    license = licenses.mit;
    maintainers = with maintainers; [ "aqtrans" ];
    platforms = platforms.linux;
  };
}
