{pkgs ? import <nixpkgs> {}}:

pkgs.mkShell {
	nativeBuildInputs = with pkgs; [
		go
		gopls
		delve
		go-tools
	];
}

