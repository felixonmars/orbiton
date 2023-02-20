#!/bin/sh
#
# Create release tarballs/zip-files
#
cd v2
name=orbiton
version=$(grep -i version main.go | head -1 | cut -d' ' -f4 | cut -d'"' -f1)
echo "Version $version"

export GOBUILD=( go build -mod=vendor -trimpath -ldflags "-w -s" -a -o )
export CGO_ENABLED=0

echo 'Compiling...'
echo '* Linux x86_64'
export GOOS=linux
GOARCH=amd64 "${GOBUILD[@]}" $name.linux_x86_64_static
echo '* Linux aarch64'
export GOARCH=arm64
"${GOBUILD[@]}" $name.linux_aarch64_static
export GOARCH=arm
echo '* Linux armv6'
GOARM=6 "${GOBUILD[@]}" $name.linux_armv6_static
echo '* Linux armv7 (RPI 2/3/4)'
GOARM=7 "${GOBUILD[@]}" $name.linux_armv7_static
unset GOARM
echo '* Linux riscv64'
GOARCH=riscv64 "${GOBUILD[@]}" $name.linux_riscv64_static

echo '* macOS x86_64'
export GOOS=darwin
GOARCH=amd64 "${GOBUILD[@]}" $name.macos_x86_64_static
echo '* macOS aarch64'
GOARCH=arm64 "${GOBUILD[@]}" $name.macos_aarch64_static

echo '* FreeBSD x86_64'
export GOOS=freebsd
GOARCH=amd64 "${GOBUILD[@]}" $name.freebsd_x86_64_static
echo '* FreeBSD aarch64'
GOARCH=arm64 "${GOBUILD[@]}" $name.freebsd_aarch64_static
export GOARCH=arm
echo '* FreeBSD armv6'
GOARM=6 "${GOBUILD[@]}" $name.freebsd_armv6_static
echo '* FreeBSD armv7'
GOARM=7 "${GOBUILD[@]}" $name.freebsd_armv7_static

echo '* NetBSD x86_64'
export GOOS=netbsd
GOARCH=amd64 "${GOBUILD[@]}" $name.netbsd_x86_64_static
echo '* NetBSD aarch64'
GOARCH=arm64 "${GOBUILD[@]}" $name.netbsd_aarch64_static
export GOARCH=arm
echo '* NetBSD armv6'
GOARM=6 "${GOBUILD[@]}" $name.netbsd_armv6_static
echo '* NetBSD armv7'
GOARM=7 "${GOBUILD[@]}" $name.netbsd_armv7_static

# OpenBSD (and Plan9) did not compile: https://github.com/pkg/term/issues/27

# Compress the Linux releases with xz
for p in \
  linux_x86_64_static linux_aarch64_static linux_armv6_static \
  linux_armv7_static linux_riscv64_static
do
  echo "Compressing $name-$version.$p.tar.xz"
  mkdir "$name-$version-$p"
  cp ../$name.1 "$name-$version-$p/"
  gzip "$name-$version-$p/$name.1"
  cp $name.$p "$name-$version-$p/$name"
  cp ../LICENSE "$name-$version-$p/"
  tar Jcf "$name-$version-$p.tar.xz" "$name-$version-$p/"
  rm -r "$name-$version-$p"
  rm $name.$p
done

# Compress the other tarballs with gz
for p in \
  macos_x86_64_static macos_aarch64_static \
  freebsd_x86_64_static freebsd_aarch64_static freebsd_armv6_static \
  freebsd_armv7_static netbsd_x86_64_static netbsd_aarch64_static \
  netbsd_armv6_static netbsd_armv7_static
do
  echo "Compressing $name-$version.$p.tar.gz"
  mkdir "$name-$version-$p"
  cp ../$name.1 "$name-$version-$p/"
  gzip "$name-$version-$p/$name.1"
  cp $name.$p "$name-$version-$p/$name"
  cp ../LICENSE "$name-$version-$p/"
  tar zcf "$name-$version-$p.tar.gz" "$name-$version-$p/"
  rm -r "$name-$version-$p"
  rm $name.$p
done
cd ..

mkdir -p release
mv -v v2/$name-$version* release
