run:
	air

clean:
	rm -rf work
	rm -rf tmp
	rm -f result

publish:
	nix build --no-link .#
	nix build --no-link .#beldum-web
	publish-version

.PHONY: run clean publish
