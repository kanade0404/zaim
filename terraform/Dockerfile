FROM hashicorp/terraform:1.6.5 AS terraform

FROM ubuntu:22.04
COPY ./aqua.yaml .
RUN apt-get update && \
    apt-get install -y curl apt-transport-https ca-certificates gnupg python3 && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update && \
    apt-get install -y google-cloud-sdk curl vim && \
    curl -sSfL -O https://raw.githubusercontent.com/aquaproj/aqua-installer/v2.2.0/aqua-installer && \
    echo "d13118c3172d90ffa6be205344b93e8621de9bf47c852d80da188ffa6985c276  aqua-installer" | sha256sum -c && \
    chmod +x aqua-installer && \
    ./aqua-installer && \
    export AQUA_ROOT_DIR="${AQUA_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/aquaproj-aqua}" && \
    export PATH="${AQUA_ROOT_DIR}/bin:$PATH" && \
    aqua i
ENV PATH="/root/.local/share/aquaproj-aqua/bin:${PATH}"
COPY --from=terraform /bin/terraform /bin/terraform
CMD ["bash"]
