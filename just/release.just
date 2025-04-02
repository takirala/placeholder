BUILD_DIR := `mktemp -d`

push-app app_name app_version registry_url:
    # Prepare file system
    mkdir -p "{{BUILD_DIR}}/services/{{app_name}}/{{app_version}}/"
    cp -r "{{justfile_directory()}}/services/{{app_name}}/{{app_version}}" "{{BUILD_DIR}}/services/{{app_name}}/"
    cd "{{BUILD_DIR}}"

    # Push artifact
    just print-ocirepository {{ app_name }} {{ app_version }} {{ registry_url }} > generated_oci_repository.yaml
    oras push {{ registry_url }}/{{ app_name }}:{{ app_version }} \
      --artifact-type application/vnd.ocirepository.catalog.nkp.ui.io.v1+yaml \
      services \
      generated_oci_repository.yaml:application/vnd.ocirepository.catalog.nkp.ui.io.v1+yaml \
      --format json
