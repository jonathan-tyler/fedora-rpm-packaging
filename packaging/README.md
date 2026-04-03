# Packaging Notes

Place package-specific inputs here.

Suggested layout:

```text
packaging/
  sesh/
    sesh.spec
    SOURCES/
  television/
    television.spec
    SOURCES/
  yazi/
    yazi.spec
    SOURCES/
```

Keep the connected preparation steps separate from the mock rebuild step:

- create source tarballs from a pinned upstream commit or tag
- vendor dependencies into a separate tarball when the package ecosystem needs it
- make the SRPM self-contained before calling mock
