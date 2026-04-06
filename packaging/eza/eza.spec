Name:           eza
Version:        0.23.4
Release:        1%{?dist}
Summary:        Modern replacement for ls
License:        EUPL-1.2
URL:            https://github.com/eza-community/eza
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}-vendor.tar.gz
BuildRequires:  cargo
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  libgit2-devel
BuildRequires:  make
BuildRequires:  openssl-devel
BuildRequires:  pandoc
BuildRequires:  perl-core
BuildRequires:  pkgconf-pkg-config
BuildRequires:  rust

%description
Eza is a modern replacement for ls with colors, Git integration, and tree views.

%prep
%autosetup -n %{name}-%{version}
tar -xzf %{SOURCE1}
mkdir -p .cargo
cp cargo-config.toml .cargo/config.toml

%build
mkdir -p target/man
for page in eza.1 eza_colors.5 eza_colors-explanation.5; do
  sed "s/\\$version/v%{version}/g" "man/${page}.md" | pandoc --standalone -f markdown -t man > "target/man/${page}"
done
export CARGO_NET_OFFLINE=true
cargo build --locked --offline --release --bin eza

%install
install -Dpm0755 target/release/eza %{buildroot}%{_bindir}/eza
install -Dpm0644 target/man/eza.1 %{buildroot}%{_mandir}/man1/eza.1
install -Dpm0644 target/man/eza_colors.5 %{buildroot}%{_mandir}/man5/eza_colors.5
install -Dpm0644 target/man/eza_colors-explanation.5 %{buildroot}%{_mandir}/man5/eza_colors-explanation.5
install -Dpm0644 completions/bash/eza %{buildroot}%{_datadir}/bash-completion/completions/eza
install -Dpm0644 completions/fish/eza.fish %{buildroot}%{_datadir}/fish/vendor_completions.d/eza.fish
install -Dpm0644 completions/zsh/_eza %{buildroot}%{_datadir}/zsh/site-functions/_eza

%files
%license LICENSE.txt LICENSES
%doc README.md CHANGELOG.md
%{_bindir}/eza
%{_mandir}/man1/eza.1*
%{_mandir}/man5/eza_colors.5*
%{_mandir}/man5/eza_colors-explanation.5*
%{_datadir}/bash-completion/completions/eza
%{_datadir}/fish/vendor_completions.d/eza.fish
%{_datadir}/zsh/site-functions/_eza

%changelog
* Sun Apr 05 2026 GitHub Copilot <copilot@example.invalid> - 0.23.4-1
- Initial local package recipe scaffold