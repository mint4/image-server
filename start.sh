#!/usr/bin/env sh

DEFAULT_QUALITY=${DEFAULT_QUALITY:-75}
EXTENSIONS=${EXTENSIONS:-jpg,gif,webp}
LISTEN=${LISTEN:-0.0.0.0}
LOCAL_BASE_PATH=${LOCAL_BASE_PATH:-public}
OUTPUTS=${OUTPUTS:-x200.jpg}
MAXIMUM_WIDTH=${MAXIMUM_WIDTH:-1000}
NAMESPACE=${NAMESPACE:-default}
PORT=${PORT:-7000}
PROC_CONCURRENCY=${PROC_CONCURRENCY:-4}
UPL_CONCURRENCY=${UPL_CONCURRENCY:-10}
MAX_FILE_AGE=${MAX_FILE_AGE:-30}
CLEAR_ORIGINALS=${CLEAR_ORIGINALS:-false}

exec bin/image-server server \
	--max_file_age          ${MAX_FILE_AGE} \
	--default_quality       ${DEFAULT_QUALITY} \
	--extensions            ${EXTENSIONS} \
	--listen                ${LISTEN} \
	--local_base_path       ${LOCAL_BASE_PATH} \
	--outputs               ${OUTPUTS} \
	--maximum_width         ${MAXIMUM_WIDTH} \
	--namespace             ${NAMESPACE} \
	--port                  ${PORT} \
	--processor_concurrency ${PROC_CONCURRENCY} \
	--uploader_concurrency  ${UPL_CONCURRENCY} \
	$([ ! "$CLEAR_ORIGINALS" = "false" ] && echo "--clear_only_originals")
