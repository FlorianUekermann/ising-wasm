#!/usr/bin/env bash
set -e

GOOS=js GOARCH=wasm go build -o ./.build/main.wasm .

cat <<EOF > ./.build/main.html
<html>
<head>
    <meta charset="UTF-8">
    <title>Go WASM</title>
</head>
<body>
    <script type="module">
EOF

cat ./.build/wasm_exec.js >> ./.build/main.html

echo >> ./.build/main.html
echo -n 'let wasm_base64 = "' >> ./.build/main.html

base64 -w 0 ./.build/main.wasm >> ./.build/main.html

cat <<EOF >> ./.build/main.html
";
        let wasm_buffer = Uint8Array.from(atob(wasm_base64), c => c.charCodeAt(0)).buffer;
        window.onload = () => {
            const go = new Go();
            WebAssembly.instantiate(wasm_buffer, go.importObject).then((result) => { go.run(result.instance); });
        }
    </script>
</body>
</html>
EOF

cp ./.build/main.html ./main.html
