# Design: Simplify S1QL queries with power-query `any()` syntax

## Problem

The S1QL backend ([pkg/backends/s1ql/s1ql.go](../../../pkg/backends/s1ql/s1ql.go)) searches multiple fields for the same value by writing explicit `||` clauses, e.g.:

```
src.process.image.md5 in ("hash1") || tgt.file.md5 in ("hash1")
```

SentinelOne's power-query syntax supports `any(field1, field2) IN (...)` for exactly this "same value, multiple fields" pattern, producing shorter, more readable queries:

```
any(src.process.image.md5, tgt.file.md5) in ("hash1")
```

## Approach

Add a single helper method to `S1QLBackend`:

```go
func (b *S1QLBackend) generateAnyFieldQuery(fields []string, values []string) string {
    return fmt.Sprintf(`any(%s) in (%s)`, strings.Join(fields, ", "), b.formatStringList(values))
}
```

Rewrite the four affected generators to call it:

- `generateMD5Query` → `any(src.process.image.md5, tgt.file.md5) in (...)`
- `generateSHA1Query` → `any(src.process.image.sha1, tgt.file.sha1) in (...)`
- `generateSHA256Query` → `any(src.process.image.sha256, tgt.file.sha256) in (...)`
- `generateIPQuery` → `any(src.ip.address, dst.ip.address) in (...)`

`generateIPQuery` currently special-cases a single IP with `field = "x" || field2 = "x"`. That special case is removed — a single IP now goes through the same `any(...) in (...)` path as multiple IPs, so there's one code path instead of two.

`generateDomainQuery` is unchanged — it only ever queried one field (`event.dns.request`), so there's nothing to consolidate.

Everything else is untouched: `formatStringList`, `escapeString`, the `GenerateQuery`/`GenerateQueries` grouping and joining logic (which combines different hash/network categories with `||`, an orthogonal concern), error handling for empty/nil IOC sets, and the >1000-IOC warning.

## Testing

Update expected strings in [pkg/backends/s1ql/s1ql_test.go](../../../pkg/backends/s1ql/s1ql_test.go):

- `TestS1QLBackend_GenerateQuery`: all MD5/SHA1/SHA256/IP cases (single and multiple), the combined-all-IOC-types case, and the single-IP case (now unified with multi-IP).
- `TestS1QLBackend_GenerateQueries`: MD5/SHA256 cases.
- Domain-only tests and `TestS1QLBackend_escapeString` are unaffected.

## Docs

Update example query outputs to the new `any(...)` syntax in:

- [README.md](../../../README.md)
- [openspec/specs/s1ql-backend/spec.md](../../../openspec/specs/s1ql-backend/spec.md) (acceptance criteria and scenarios for MD5/SHA1/SHA256/IPv4 requirements)

This does not go through the full OpenSpec change-proposal workflow (proposal.md/tasks.md/delta specs) — by explicit user choice, this is a superpowers-brainstorming design doc instead, but the OpenSpec `specs/` content is still kept in sync since it documents current, deployed behavior.

## Out of scope

- The `r7` backend (unaffected — this is S1QL-specific power-query syntax).
- Any change to how different IOC categories (hashes vs. domains vs. IPs) are joined together in `GenerateQuery`.
