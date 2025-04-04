BUILD_DIR := `mktemp -d`

print-ocirepository-name app_name app_version:
    @echo "{{app_name}}-{{ app_version }}" | tr '.' '-'

print-kustomization app_name app_version registry_url:
    @repo_name=$(just print-ocirepository-name {{ app_name }} {{ app_version }}) && \
    flux create kustomization ${repo_name} \
      --namespace \${releaseNamespace} \
      --source OCIRepository/${repo_name} \
      --path ./services/{{app_name}}/{{app_version}} \
      --prune=true \
      --interval=1m \
      --export

print-ocirepository app_name app_version registry_url:
    @repo_name=$(just print-ocirepository-name {{ app_name }} {{ app_version }}) && \
    flux create source oci ${repo_name} \
      --namespace \${releaseNamespace} \
      --label "kommander.d2iq.io/ocirepository-type=catalog" \
      --url oci://{{registry_url}}/{{app_name}} \
      --tag {{app_version}} \
      --export

build-app app_name app_version registry_url:
    # Prepare file system
    @mkdir -p "{{BUILD_DIR}}/services/{{app_name}}/{{app_version}}/"
    @cp -r "{{justfile_directory()}}/services/{{app_name}}/{{app_version}}" "{{BUILD_DIR}}/services/{{app_name}}/"
    @export buildSubstitute_ociRepositoryName=$(just print-ocirepository-name {{ app_name }} {{ app_version }}) && \
    for file in $(find {{BUILD_DIR}}/services/{{app_name}} -type f); do \
      envsubst \${buildSubstitute_ociRepositoryName} <$file >$file.tmp; \
      mv $file.tmp $file; \
    done
    @just print-ocirepository {{ app_name }} {{ app_version }} {{ registry_url }} > {{BUILD_DIR}}/generated_oci_repository.yaml
    @echo "OCI created at {{BUILD_DIR}}/generated_oci_repository.yaml"

push-app app_name app_version registry_url: (build-app app_name app_version registry_url)
    # Push the artifact
    @cd "{{BUILD_DIR}}"
    oras push {{ registry_url }}/{{ app_name }}:{{ app_version }} \
      --artifact-type application/vnd.ocirepository.catalog.nkp.io.v1+yaml \
      services \
      generated_oci_repository.yaml:application/vnd.ocirepository.catalog.nkp.io.v1+yaml \
      --format json
