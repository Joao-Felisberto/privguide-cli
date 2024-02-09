{pkgs ? import <nixpkgs> {}}:

pkgs.mkShell {
	nativeBuildInputs = with pkgs; [
		#apache-jena-fuseki

		go

		gopls
		delve
		go-tools
	];
}

