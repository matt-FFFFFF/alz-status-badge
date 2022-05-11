# alz-status-badge

Application to report validity of ALZ variants using shields.io badges.

HTTP API on path `/api/badge`

## Data

API reads data in JSON format from this repo in: `/data/approved-variants.json`.

Data is a JSON array, e.g.

```json
[
  "variant1",
  "variant2"
]
```

The data is refreshed periodically, see [`ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL`](#ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL)

If the supplied variant is in the list, the response will be a green badge. If not, then the response will be a red badge with unapproved status.

## Usage

* Make A `GET` request and add a `variant` query parameter to the URL.
* The response will be a SVG badge with the status of the variant.

### Environment Variables

The application can be configured using environment variables:

#### `ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL`

A time duration to refresh the approved variants data in memory. Default is `15m`.

A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

#### `ALZSTATUSBADGE_APPROVED_VARIANTS_URL`

URL to the JSON file with the approved variants. Default is this repo's `/data/approved-variants.json` file.

#### `ALZSTATUSBADGE_LISTEN_ADDRESS`

Listen address for HTTP server. Default is `:8080`.

## Example embedding in markdown

```markdown
![badge](https://your.url/api/badge/?variant=myvariant)
```
