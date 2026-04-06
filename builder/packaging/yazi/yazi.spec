Name:           yazi
Version:        26.1.22
Release:        1%{?dist}
Summary:        Blazing fast terminal file manager written in Rust
License:        MIT
URL:            https://github.com/sxyazi/yazi
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}-vendor.tar.gz
BuildRequires:  cargo
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  make
BuildRequires:  oniguruma-devel
BuildRequires:  perl-core
BuildRequires:  pkgconf-pkg-config
BuildRequires:  rust

%description
Yazi is a terminal file manager written in Rust with non-blocking async I/O.

%prep
%autosetup -n %{name}-%{version}
tar -xzf %{SOURCE1}
mkdir -p .cargo
cp cargo-config.toml .cargo/config.toml

%build
export CARGO_BUILD_JOBS=1
export CARGO_NET_OFFLINE=true
cargo build --locked --offline --release --bin yazi --bin ya

%install
install -Dpm0755 target/release/yazi %{buildroot}%{_bindir}/yazi
install -Dpm0755 target/release/ya %{buildroot}%{_bindir}/ya

%files
%license LICENSE LICENSE-ICONS
%doc README.md CHANGELOG.md
%{_bindir}/yazi
%{_bindir}/ya

%changelog
* Sat Apr 04 2026 GitHub Copilot <copilot@example.invalid> - 26.1.22-1
- Initial local package recipe scaffold