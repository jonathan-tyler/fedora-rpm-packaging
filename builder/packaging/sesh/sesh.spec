Name:           sesh
Version:        2.24.2
Release:        1%{?dist}
Summary:        Smart tmux session manager
License:        MIT
URL:            https://github.com/joshmedeski/sesh
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}-vendor.tar.gz
BuildRequires:  golang

%description
Sesh is a CLI that helps create and manage tmux sessions quickly and easily
using zoxide.

%prep
%autosetup -n %{name}-%{version}
tar -xzf %{SOURCE1}

%build
export CGO_ENABLED=0
go build -mod=vendor -buildvcs=false -ldflags "-X main.version=v%{version}" -o sesh ./

%install
install -Dpm0755 sesh %{buildroot}%{_bindir}/sesh

%files
%license LICENSE
%doc README.md
%{_bindir}/sesh

%changelog
* Sat Apr 04 2026 GitHub Copilot <copilot@example.invalid> - 2.24.2-1
- Initial local package recipe scaffold