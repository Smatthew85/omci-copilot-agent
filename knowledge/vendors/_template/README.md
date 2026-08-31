# How to Add a Vendor Quirk

This folder contains the contribution template for documenting a vendor-specific OMCI
deviation.

---

## Steps

1. **Copy** `_template/quirk.md` into the appropriate vendor folder
   (e.g., `knowledge/vendors/nokia/`).

2. **Rename** the copy to `NNN-short-slug.md` where `NNN` is a zero-padded sequence
   number unique within that vendor folder.
   Example: `001-me-171-truncates-rules-to-16-bytes.md`

3. **Fill in every section.** Do not leave any section blank. If a field genuinely does
   not apply, write "N/A" and explain why.

4. **Do not omit sources.** Every quirk entry requires at least one verifiable source
   (vendor bug ID, VOLTHA issue URL, internal ticket, reproducible test case, or a dated
   field observation that includes device model and firmware version).

5. **Open a PR** targeting `main`. Reference the source ticket or issue in the PR
   description.

---

## Naming Convention

```
knowledge/vendors/<vendor>/NNN-short-slug.md
```

- `<vendor>` — one of `adtran`, `nokia`, `calix`, `huawei`, `zte`, or `_other`
- `NNN` — three-digit zero-padded sequence per vendor (start at `001`)
- `short-slug` — lowercase, hyphenated, ≤ 60 chars, describes the quirk
  (e.g., `me-171-truncates-rules-to-16-bytes`, `tcont-create-returns-0x03`)

---

## Template File

[`quirk.md`](./quirk.md)
