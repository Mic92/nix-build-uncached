# nix-build-uncached

![Test](https://github.com/Mic92/nix-build-uncached/workflows/Test/badge.svg)

`nix-build` by default will download already built packages, resulting in
unnecessary downloads even if no package has been changed.
`nix-build-uncached` will only build packages not yet in binary caches.

## USAGE

`nix-build-uncached` is available in nixpkgs.

In the following example ci.nix contains all expressions
that should be built. Since only `hello-nur` is not yet in
the binary cache, all other packages are skipped.

```
[joerg@turingmachine] nix-build-uncached ci.nix
$ nix-build --dry-run ci.nix --keep-going
these derivations will be built:
  /nix/store/s5alllpjx9fmdj26mf9cmxzs3xyxjn7f-hello-2.00.tar.gz.drv
  /nix/store/03m2lwg4zia58zqm7hqlb3r0cgfq53cn-hello-2.00.drv
these paths will be fetched (198.28 MiB download, 681.37 MiB unpacked):
  /nix/store/0mijq2b50xmgk6akxdrg3x8x0k7784jb-python3.8-kiwisolver-1.2.0
  /nix/store/0q1g21q160w2cj2745r24pfn2yb8pmda-python3.8-jupyter_client-6.1.5
  /nix/store/0qxi1gzv1xgpxy58nfk9pxlfqybv1198-cntr-1.2.0
 # ...
$ nix build --keep-going /nix/store/s5alllpjx9fmdj26mf9cmxzs3xyxjn7f-hello-2.00.tar.gz.drv /nix/store/03m2lwg4zia58zqm7hqlb3r0cgfq53cn-hello-2.00.drv
[1/2 built, 1 copied (0.7 MiB)] connecting to 'ssh://Mic92@prism.r'

```

We can pass all arguments that also `nix-build` accept:

```
[joerg@turingmachine] nix-build-uncached -E 'with import <nixpkgs> {}; hello'
$ nix-build --dry-run -E with import <nixpkgs> {}; hello --keep-going
these paths will be fetched (0.04 MiB download, 0.20 MiB unpacked):
  /nix/store/aldyr0pjzqydf1vn9lzz7p5gvc141fhn-hello-2.10
```

However this only affects the evaluation during the dry build, if you want to
pass arguments to the final `nix build` instead, use `-build-flags`:

```
[joerg@turingmachine] nix-build-uncached -build-flags '--builders ""' ci.nix
$ nix-build --dry-run ci.nix --builders
these derivations will be built:
  /nix/store/s5alllpjx9fmdj26mf9cmxzs3xyxjn7f-hello-2.00.tar.gz.drv
  /nix/store/03m2lwg4zia58zqm7hqlb3r0cgfq53cn-hello-2.00.drv
these paths will be fetched (198.28 MiB download, 681.37 MiB unpacked):
  /nix/store/0mijq2b50xmgk6akxdrg3x8x0k7784jb-python3.8-kiwisolver-1.2.0
  /nix/store/0q1g21q160w2cj2745r24pfn2yb8pmda-python3.8-jupyter_client-6.1.5
  /nix/store/0qxi1gzv1xgpxy58nfk9pxlfqybv1198-cntr-1.2.0
 # ...
$ nix build --builders   /nix/store/s5alllpjx9fmdj26mf9cmxzs3xyxjn7f-hello-2.00.tar.gz.drv /nix/store/03m2lwg4zia58zqm7hqlb3r0cgfq53cn-hello-2.00.drv
[1/2 built, 1 copied (0.7 MiB)] connecting to 'ssh://Mic92@prism.r'
```

### Flakes

We cannot support flakes directly at the time because `nix-build` does
not accept those and `nix build`'s `--no-dry-run` is broken.
However it is possible to add a wrapper nix expression that imports a flake.
The following example imports hydra jobs from a nix flake the same directory.
It assumes that the flake also has `nixpkgs` in its inputs.

```nix
let
  outputs = builtins.getFlake (toString ./.);
  pkgs = outputs.inputs.nixpkgs;
  drvs = pkgs.lib.collect pkgs.lib.isDerivation outputs.hydraJobs;
in drvs
```

```console
$ nix-build-uncached ./ci.nix
```

### Packages with `allowSubstitutes = false`

If your package set you are building has packages at the top level scope that
have the attribute `allowSubstitutes = false;` set, than `nix-build-uncached`
will build/download them everytime. This attribute is set for some builders such
as `writeText` or `writeScriptBin`. A workaround is to use the following
[nix library](./scripts/force_cached.nix) and save it as
`force_cached.nix`. Than wrap your attribute set like this:

```nix
let
  pkgs = import <nixpkgs> {};
in (pkgs.callPackage ./force_cached.nix {}) {
  hello = pkgs.writeScriptBin "hello" ''
    #!/bin/sh
    exec ${pkgs.hello}/bin/hello
  '';
}
```


## Real-world examples

- [Using nix-build-uncached in github actions](https://github.com/Mic92/nur-packages/blob/master/.github/workflows/build.yml)
