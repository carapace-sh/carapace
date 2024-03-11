export SHELL=bash
export STARSHIP_SHELL=bash
export LS_COLORS="$(vivid generate dracula)"
[[ ! -z $BLE ]] && source /opt/ble.sh/out/ble.sh 
eval "$(starship init bash)"
[[ -z "${TARGET}" ]] || source <(${TARGET} _carapace)
