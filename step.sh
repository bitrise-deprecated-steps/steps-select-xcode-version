#!/bin/bash

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# load bash utils
source "${THIS_SCRIPT_DIR}/bash_utils/utils.sh"
source "${THIS_SCRIPT_DIR}/bash_utils/formatted_output.sh"


# ------------------------------
# --- Error Cleanup

function finalcleanup {
  echo "-> finalcleanup"
  local fail_msg="$1"

  write_section_to_formatted_output "# Error"
  if [ ! -z "${fail_msg}" ] ; then
    write_section_to_formatted_output "**Error Description**:"
    write_section_to_formatted_output "${fail_msg}"
  fi
  write_section_to_formatted_output "*See the logs for more information*"
}

function CLEANUP_ON_ERROR_FN {
  local err_msg="$1"
  finalcleanup "${err_msg}"
}
set_error_cleanup_function CLEANUP_ON_ERROR_FN


# ------------------------------
# --- Main

function set_xcode_path_by_channel {
    local channel_id="$1"
    local channel_version_map_file_pth="/Applications/Xcodes/version-map-${channel_id}"
    mapping_file_content="$(cat ${channel_version_map_file_pth})"
    if [[ "${mapping_file_content}" == "" ]] ; then
      echo " [!] No version mapping found for channel: ${channel_id} / mapping file path: ${channel_version_map_file_pth}"
      exit 2
    fi
    CONFIG_xcode_path="${mapping_file_content}"
}

if [ -z "${version_channel_id}" ] ; then
	finalcleanup "No Xcode Version Channel ID specified!"
	exit 1
fi

write_section_to_formatted_output "# Select Xcode Version"

CONFIG_xcode_path=""
set_xcode_path_by_channel "${version_channel_id}"

echo_string_to_formatted_output " * Selecting Xcode: \`${CONFIG_xcode_path}\`"


sudo xcode-select --switch "${CONFIG_xcode_path}"
fail_if_cmd_error "Failed to activate the specified Xcode version"

# --- Report
xcode_version_info_text="$(xcodebuild -version)"
indented_xcode_version_info_text=$(printf %s "${xcode_version_info_text}" | awk '{print "    " $0}')
xcode_sdks_info_text="$(xcodebuild -showsdks)"
indented_xcode_sdks_info_text=$(printf %s "${xcode_sdks_info_text}" | awk '{print "    " $0}')
#
write_section_to_formatted_output '# Xcode Information'
write_section_to_formatted_output '## Xcode Version'
echo_string_to_formatted_output "${indented_xcode_version_info_text}"
write_section_to_formatted_output '## Xcode SDKs'
echo_string_to_formatted_output "${indented_xcode_sdks_info_text}"
