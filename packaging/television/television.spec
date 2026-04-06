Name:           television
Version:        0.15.4
Release:        1%{?dist}
Summary:        Fast fuzzy finder for the terminal
License:        MIT
URL:            https://github.com/alexpasmantier/television
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}-vendor.tar.gz
BuildRequires:  cargo
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  make
BuildRequires:  perl-core
BuildRequires:  pkgconf-pkg-config
BuildRequires:  rust

%description
Television is a fast, portable, and hackable fuzzy finder for the terminal.

%prep
%autosetup -n %{name}-%{version}
tar -xzf %{SOURCE1}
mkdir -p .cargo
cp cargo-config.toml .cargo/config.toml

%build
export CARGO_NET_OFFLINE=true
cargo build --locked --offline --release --bin tv

%install
install -Dpm0755 target/release/tv %{buildroot}%{_bindir}/tv
install -Dpm0644 man/tv.1 %{buildroot}%{_mandir}/man1/tv.1

%files
%license LICENSE
%doc README.md CHANGELOG.md
%{_bindir}/tv
%{_mandir}/man1/tv.1*

%changelog
* Sat Apr 04 2026 GitHub Copilot <copilot@example.invalid> - 0.15.4-1
- Initial local package recipe scaffold