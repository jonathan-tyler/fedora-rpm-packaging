%global debug_package %{nil}
%global _missing_build_ids_terminate_build 0
%global zig_version 0.15.2

Name:           ghostty
Version:        1.3.1
Release:        1%{?dist}
Summary:        Fast, feature-rich terminal emulator
License:        MIT
URL:            https://github.com/ghostty-org/ghostty
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}-zig-cache.tar.gz
Source2:        zig-x86_64-linux-%{zig_version}.tar.xz
BuildRequires:  blueprint-compiler
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  gettext
BuildRequires:  gtk4-devel
BuildRequires:  gtk4-layer-shell-devel
BuildRequires:  libadwaita-devel
BuildRequires:  make
BuildRequires:  ncurses
BuildRequires:  oniguruma-devel
BuildRequires:  pandoc
BuildRequires:  pkgconf-pkg-config

%description
Ghostty is a fast, feature-rich terminal emulator with GTK integration,
desktop assets, shell integration helpers, and bundled themes.

%prep
%autosetup -n %{name}-%{version}
mkdir -p .fpb
tar -xzf %{SOURCE1} -C .fpb
tar -xJf %{SOURCE2} -C .fpb

%build
export PATH="$PWD/.fpb/zig-x86_64-linux-%{zig_version}:$PATH"
zig build \
	--system "$PWD/.fpb/p" \
	-fno-sys=fontconfig \
	-fno-sys=freetype \
	-fno-sys=harfbuzz \
	-fno-sys=libpng \
	-fno-sys=libxml2 \
	-fno-sys=oniguruma \
	-fno-sys=zlib \
	-Doptimize=ReleaseFast \
	-Dcpu=baseline \
	-Dstrip=false

%install
export PATH="$PWD/.fpb/zig-x86_64-linux-%{zig_version}:$PATH"
DESTDIR=%{buildroot} \
zig build \
	--prefix %{_prefix} \
	--system "$PWD/.fpb/p" \
	-fno-sys=fontconfig \
	-fno-sys=freetype \
	-fno-sys=harfbuzz \
	-fno-sys=libpng \
	-fno-sys=libxml2 \
	-fno-sys=oniguruma \
	-fno-sys=zlib \
	-Doptimize=ReleaseFast \
	-Dcpu=baseline \
	-Dstrip=false \
	install
install -Dpm0755 zig-out/bin/ghostty %{buildroot}%{_bindir}/ghostty

%files
%license LICENSE
%doc README.md PACKAGING.md
%{_bindir}/ghostty
%{_includedir}/ghostty
%{_prefix}/lib/libghostty-vt.so*
%{_prefix}/lib/systemd/user/app-com.mitchellh.ghostty.service
%{_datadir}/applications/com.mitchellh.ghostty.desktop
%{_datadir}/bash-completion/completions/ghostty.bash
%{_datadir}/bat/syntaxes/ghostty.sublime-syntax
%{_datadir}/dbus-1/services/com.mitchellh.ghostty.service
%{_datadir}/fish/vendor_completions.d/ghostty.fish
%{_datadir}/ghostty
%{_datadir}/icons/hicolor/*/apps/com.mitchellh.ghostty.png
%{_datadir}/kio/servicemenus/com.mitchellh.ghostty.desktop
%{_datadir}/locale/*/LC_MESSAGES/com.mitchellh.ghostty.mo
%{_datadir}/metainfo/com.mitchellh.ghostty.metainfo.xml
%{_datadir}/nautilus-python/extensions/ghostty.py
%{_datadir}/nvim/site
%{_datadir}/pkgconfig/libghostty-vt.pc
%{_datadir}/terminfo/g/ghostty
%{_datadir}/terminfo/x/xterm-ghostty
%{_datadir}/vim/vimfiles
%{_datadir}/zsh/site-functions/_ghostty
%{_mandir}/man1/ghostty.1*
%{_mandir}/man5/ghostty.5*

%changelog
* Sat Apr 04 2026 GitHub Copilot <copilot@example.invalid> - 1.3.1-1
- Initial local package recipe scaffold