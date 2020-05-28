# nix-build-uncached

![Test](https://github.com/Mic92/nix-build-uncached/workflows/Test/badge.svg)

`nix-build` by default will download already built packages, resulting in
unnecessary downloads even if no package has been changed.
`nix-build-uncached` will only build packages not yet in binary caches.

## USAGE

`nix-build-uncached` is available in nixpkgs-unstable.

```
$ nix-shell -p nix-build-uncached --run 'nix-build-uncached --help'
Usage of nix-build-uncached:
  -build-flags string
    	additional arguments to pass to both nix build (default "--keep-going")
  -flags string
    	additional arguments to pass to both nix-env/nix build
```

Pass a file with the nix expressions you want to build.
As a result `nix-build-uncached` will build all packages,
not present in the binary cache:

```console
[joerg@turingmachine] nix-build-uncached ./ci.nix
$ nix-env -f non-broken.nix --drv-path -qaP * --xml --meta
$ nix-build --dry-run non-broken.nix
1/40 attribute(s) will be built:
  hello-nur
$ nix build -f /tmp/859272287.nix --keep-going
[1 built, 1 copied (0.2 MiB)]
```

## Real-world examples

- [Using nix-build-uncached in github actions](https://github.com/Mic92/nur-packages/blob/master/.github/workflows/build.yml)
