rfc1123_version app_version:
    @echo "{{ app_version }}" | tr '.' '-'

print-kustomization app_name app_version registry_url:
    #!/usr/bin/env bash
    set -euo pipefail

    rfc1123_version=$(just rfc1123_version {{ app_version }})
    flux create kustomization {{app_name}}-${rfc1123_version} \
      --namespace \${releaseNamespace} \
      --source OCIRepository/{{app_name}}-${rfc1123_version} \
      --path ./services/{{app_name}}/{{app_version}} \
      --prune=true \
      --interval=1m \
      --export

print-ocirepository app_name app_version registry_url:
    #!/usr/bin/env bash
    set -euo pipefail

    rfc1123_version=$(just rfc1123_version {{ app_version }})
    flux create source oci {{app_name}}-${rfc1123_version} \
      --namespace \${releaseNamespace} \
      --label "kommander.d2iq.io/ocirepository-type=catalog" \
      --url oci://{{registry_url}}/{{app_name}} \
      --tag {{app_version}} \
      --export
