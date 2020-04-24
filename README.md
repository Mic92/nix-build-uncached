# nix-build-uncached

`nix-build` by default will download already built packages, resulting in
unnecessary downloads even if no package has been changed.
Nix-build-uncached will only build packages not yet in binary caches.

## USAGE

Pass a file with the nix expressions you want to build.

```console
$ nix-build-uncached ./ci.nix
```

As a result `nix-build-uncached` will built all packages,
not present in the binary cache.

```
$ nix-env -f ./ci.nix --drv-path -qaP * --xml --meta
$ nix-build --dry-run ./ci.nix.nix
/nix/store/j68qphj95cx65xsxadgmy9wa08dlbjrq-hello-2.10.drv
1/40 attribute(s) will be built:
  hello-nur
$ nix-build ./ci.nix -k -A hello-nur
these derivations will be built:
  /nix/store/j68qphj95cx65xsxadgmy9wa08dlbjrq-hello-2.10.drv
building '/nix/store/j68qphj95cx65xsxadgmy9wa08dlbjrq-hello-2.10.drv' on 'ssh://nix@martha.r'...
copying path '/nix/store/3x7dwzq014bblazs7kq20p9hyzz0qh8g-hello-2.10.tar.gz' from 'https://cache.nixos.org'...
unpacking sources
unpacking source archive /nix/store/3x7dwzq014bblazs7kq20p9hyzz0qh8g-hello-2.10.tar.gz
source root is hello-2.10
setting SOURCE_DATE_EPOCH to timestamp 1416139241 of file hello-2.10/ChangeLog
patching sources
configuring
...
shrinking RPATHs of ELF executables and libraries in /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10
shrinking /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10/bin/hello
gzipping man pages under /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10/share/man/
strip is /nix/store/p1y0xl8dp4s1x1vvxxb5sn84wj6lsh8s-binutils-2.31.1/bin/strip
stripping (with command strip and flags -S) in /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10/bin
patching script interpreter paths in /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10
checking for references to /build/ in /nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10...
/nix/store/k2z0nxgz2pm2d6lbc6v7r0gfxvil5ns2-hello-2.10
```
