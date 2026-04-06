Name:           lazygit
Version:        0.60.0
Release:        1%{?dist}
Summary:        Simple terminal UI for Git commands
License:        MIT
URL:            https://github.com/jesseduffield/lazygit
Source0:        %{name}-%{version}.tar.gz
BuildRequires:  golang

%description
Lazygit provides a terminal user interface for common Git workflows.

%prep
%autosetup -n %{name}-%{version}

%build
export CGO_ENABLED=0
go build -mod=vendor -buildvcs=false -ldflags "-X main.version=v%{version}" -o lazygit ./

%install
install -Dpm0755 lazygit %{buildroot}%{_bindir}/lazygit

%files
%license LICENSE
%doc README.md
%{_bindir}/lazygit

%changelog
* Sun Apr 05 2026 GitHub Copilot <copilot@example.invalid> - 0.60.0-1
- Initial local package recipe scaffold