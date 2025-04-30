### 2. Unbounded HTTP Response Size (Denial-of-Service)
- **Description**: The fetcher reads entire HTTP response bodies into memory via `io.ReadAll`, with no configurable size limit. A very large response could exhaust memory.
- **Severity**: Medium (local CLI; attacker-controlled feed could trigger OOM).
- **Mitigation**:
  - Wrap response bodies in an `io.LimitedReader` (e.g. configurable max bytes).
  - Fail gracefully if responses exceed a reasonable size.


### 4. HTML Content Sanitization (Cross-Site Scripting)
- **Description**: The pipeline does not sanitize article HTML content. The EPUB formatter inserts `article.Content` via `template.HTML`, and the Markdown converter may pass through HTML blocks. Malicious `<script>` or `<iframe>` tags will end up in the generated output.
- **Severity**: Medium
- **Mitigation**:
  - Integrate an HTML sanitizer (e.g. [bluemonday]) before conversion, stripping unsafe tags and attributes.
  - Explicitly remove `<script>`, `<iframe>`, and event handler attributes.
