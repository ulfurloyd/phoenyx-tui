{
  lib,
  buildGoModule,
  go_1_27,
  src,
}:
(buildGoModule.override {
  go = go_1_27;
})
  {
    pname = "phoenyx-tui";
    version = "0.1.0";
    inherit src;
    subPackages = [ "cmd/phoenyx" ];
    vendorHash = "sha256-bxnj5fGcXAp1veo600mmq5f7SM2D4/WNekfWFYu8jOk=";
    doCheck = true;
    meta = {
      description = "A terminal interface for managing phoenyx systems and projects";
      mainProgram = "phoenyx";
      platforms = lib.platforms.linux;
    };
  }
