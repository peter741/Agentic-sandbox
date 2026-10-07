# Security reporting and limitations

This project runs agent-controlled commands in containers. Containers are not an
absolute security boundary. Mounted files, forwarded credentials, published
ports, and access to the Docker daemon affect what an agent can reach.

`agbox doctor` checks selected setup and connectivity conditions. It does not
audit escape resistance, firewall rules, credential isolation, image integrity,
or the permissions of the running daemon.

Do not place secrets or exploit details in public issues. Use GitHub private
vulnerability reporting if enabled. If it is not enabled, use a private reporting
channel published by the maintainer. No dedicated private contact has been
published for this derivative yet; upstream arrangements do not automatically
cover it.

Before wider distribution, enable private vulnerability reporting and document
a monitored contact. No supported-release policy or response-time guarantee is
currently published for this derivative.
