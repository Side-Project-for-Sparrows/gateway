#!/usr/bin/env bash
set -Eeuo pipefail

# ---------- Defaults ----------
: "${OUT_DIR:=.}"        # output dir
: "${DAYS:=3650}"               # 10 years
: "${CA_CN:=Local Dev CA}"      # local CA CN
: "${CN:=localhost}"            # server cert CN
: "${CURVE:=prime256v1}"        # EC curve
: "${FORCE:=0}"                 # 1 to overwrite existing files
# ------------------------------

# Ensure UTF-8 locale (best-effort; harmless if already set)
export LC_ALL="${LC_ALL:-C.UTF-8}"
export LANG="${LANG:-C.UTF-8}"

# if [ "$(id -u)" -ne 0 ]; then
#   echo "[-] Need root. Run with sudo."
#   exit 1
# fi

command -v openssl >/dev/null 2>&1 || { echo "[-] openssl not found"; exit 1; }

mkdir -p "${OUT_DIR}"
chmod 777 "${OUT_DIR}"
umask 077

CA_KEY="${OUT_DIR}/ca.key"
CA_CRT="${OUT_DIR}/ca.crt"
CA_SRL="${OUT_DIR}/ca.srl"
TLS_KEY="${OUT_DIR}/tls.key"
TLS_CSR="${OUT_DIR}/tls.csr"
TLS_CRT="${OUT_DIR}/tls.crt"
SAN_EXT="${OUT_DIR}/san.ext"

# maybe_overwrite() {
#   local f="$1"
#   if [ -e "$f" ] && [ "${FORCE}" != "1" ]; then
#     echo "[-] Exists: $f (set FORCE=1 to overwrite)"
#     exit 1
#   fi
# }
# for f in "$CA_KEY" "$CA_CRT" "$CA_SRL" "$TLS_KEY" "$TLS_CSR" "$TLS_CRT" "$SAN_EXT"; do
#   maybe_overwrite "$f"
# done

echo "[*] Generate local CA: CN='${CA_CN}', days=${DAYS}"
openssl ecparam -name "${CURVE}" -genkey -noout -out "${CA_KEY}"
chmod 777 "${CA_KEY}"
openssl req -x509 -new -key "${CA_KEY}" -sha256 -days "${DAYS}" \
  -subj "/CN=${CA_CN}" -out "${CA_CRT}"
chmod 777 "${CA_CRT}"

echo "[*] Generate server key: CN='${CN}'"
openssl ecparam -name "${CURVE}" -genkey -noout -out "${TLS_KEY}"
chmod 777 "${TLS_KEY}"

echo "[*] Generate CSR"
if ! openssl req -new -key "${TLS_KEY}" -subj "/CN=${CN}" \
     -addext "subjectAltName=DNS:${CN},DNS:localhost,IP:127.0.0.1,IP:::1" \
     -out "${TLS_CSR}" 2>/dev/null; then
  echo "[!] OpenSSL -addext not supported; injecting SAN at sign step"
  openssl req -new -key "${TLS_KEY}" -subj "/CN=${CN}" -out "${TLS_CSR}"
fi

cat > "${SAN_EXT}" <<EOF
subjectAltName=DNS:${CN},DNS:localhost,IP:127.0.0.1,IP:::1
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
basicConstraints=CA:FALSE
EOF

echo "[*] Sign server cert with local CA: days=${DAYS}"
openssl x509 -req -in "${TLS_CSR}" -CA "${CA_CRT}" -CAkey "${CA_KEY}" \
  -CAcreateserial -CAserial "${CA_SRL}" \
  -out "${TLS_CRT}" -days "${DAYS}" -sha256 -extfile "${SAN_EXT}"
chmod 777 "${TLS_CRT}"

printf "\n[OK] Wrote files to %s\n" "${OUT_DIR}"
ls -l "${CA_CRT}" "${TLS_CRT}" "${TLS_KEY}"

printf "\n[i] Cert info:\n"
openssl x509 -in "${TLS_CRT}" -noout -subject -issuer -dates
