# TransparnSee 🔍

A very simple certificate transparency log scraper for a very specific use case. I wanted to prove that it's easy to find out all other domains that belong to the same Cloudflare account.

## How does it work?

The idea is that most Cloudflare accounts get two nameservers from a large pool and those don't change often. Therefore the likeliness of two domains belonging to the same account is high if they share the same nameservers. This tool tries to find out all domains that share the same nameservers as the one you specify.

1. It scrapes through the entire certificate transparency log it was given and parses all leaf certificates using the help of the amazing [certificate-transparency-go](https://github.com/google/certificate-transparency-go).
2. For convenience it narrows further search down to first level domains (ie. `cloudflare.com`, but not `whatever.cloudflare.com`) and specific domain ending specified in config (ie. `.hu`). Finally it filters out certificate authorities that are not used by Cloudflare. I used [this](https://developers.cloudflare.com/ssl/reference/certificate-authorities/) as a reference.
3. For all domains left to try, it uses Cloudflare's DoH DNS service to determine the two nameserver records belonging to the domain. If it's a match to both nameservers set in the config, it's a match.
4. Found domains are written to a file.

## Usage

Docker/Podman is recommended.

1. Alter [config.json](config/config.json) to your liking.
2. Build the image: `docker build -t transparnsee:latest .`
3. Run the image: `docker run -v $(pwd)/output:/app/output transparnsee:latest`

For more details on this, see my article in Hungarian [here](https://logout.hu/bejegyzes/a_tls_mint_osint_nyilt_hirszerzo_eszkoz.html).