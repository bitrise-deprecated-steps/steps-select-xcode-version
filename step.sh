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
    echo " * Specified channel: ${channel_id}"
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

canonical_xcode_path="/Applications/Xcode.app"
if [[ "$CONFIG_xcode_path" != "$canonical_xcode_path" ]] ; then
  if [[ -L "${canonical_xcode_path}" ]]; then
    echo " (i) Symlink found at ${canonical_xcode_path} - removing it to replace with the selected version's symlink"
    echo '$' rm "${canonical_xcode_path}"
    rm "${canonical_xcode_path}"
  fi

  if [[ -d "${canonical_xcode_path}" ]]; then
    echo " [!] Xcode already installed at ${canonical_xcode_path}, can't create symlink for selected version!"
    exit 1
  fi

  echo " (i) Creating a symlink to the canonical path ($canonical_xcode_path), pointing to the selected version ($CONFIG_xcode_path)"
  echo '$' ln -s "${CONFIG_xcode_path}" "${canonical_xcode_path}"
  ln -s "${CONFIG_xcode_path}" "${canonical_xcode_path}"
else
  echo " (i) Selected Xcode is already at ${canonical_xcode_path} - no symlink required"
fi

echo '$' sudo xcode-select --switch "${CONFIG_xcode_path}"
sudo xcode-select --switch "${CONFIG_xcode_path}"
fail_if_cmd_error "Failed to activate the specified Xcode version"

# --- Report
#
echo
echo
echo '# Selected Xcode Information'
echo
echo
echo '## Xcode Version'
echo
xcodebuild -version
echo
echo
echo '## Xcode SDKs'
xcodebuild -showsdks
echo
echo
echo '## Removing duplicate Simulator definitions...'
echo
go run "${THIS_SCRIPT_DIR}/simctl-cleanup/main.go" -delete
echo
echo '## Simulators'
echo '(List of Simulators available for this Xcode version)'
echo
xcrun simctl list | grep -i --invert-match 'unavailable'
