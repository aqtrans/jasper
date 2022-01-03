{ config, lib, pkgs, ... }:

with lib;

let 
  cfg = config.services.jasper;
  pkgDesc = "A meme as a service.";

in {

  options = {
    services.jasper = {

      enable = mkEnableOption "${pkgName}";

      user = mkOption {
        type = types.str;
        default = "jasper";
        description = "gorram user";
      };

      group = mkOption {
        type = types.str;
        default = "jasper";
        description = "gorram group";
      };

      listenAddress = mkOption {
        type = types.str;
        default = "127.0.0.1:8002";
        description = "Listen address, IP and port";
      };    

    };
  };

  config = mkIf cfg.enable {

        users.users.${cfg.user} = {
          name = cfg.user;
          group = cfg.group;
          isSystemUser = true;
          createHome = false;
          description = pkgDesc;
        };

        users.groups.${cfg.user} = {
          name = cfg.group;
        };

        systemd.services.jasper = {
          description = pkgDesc;
          wantedBy = [ "multi-user.target" ];
          after = [ "network-online.target" ];
          serviceConfig = {
            User = cfg.user;
            Group = cfg.group;
            Restart = "always";
            ProtectSystem = "strict";
            ExecStart = ''
              ${pkgs.jasper}/bin/jasper
            '';
          };
        };

  };

}
