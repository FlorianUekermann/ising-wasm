use std::{
    fs::File,
    io::{copy, Write},
};

use base64::{engine::general_purpose::STANDARD, write::EncoderWriter};

fn main() {
    let mut wasm = File::open("ising_bg.wasm").unwrap();
    let mut js = File::open("ising.js").unwrap();
    let mut html = File::create("main.html").unwrap();
    html.write_all(HTML_HEAD.as_bytes()).unwrap();
    copy(&mut js, &mut html).unwrap();
    html.write_all(HTML_MID.as_bytes()).unwrap();
    let mut enc = EncoderWriter::new(&mut html, &STANDARD);
    copy(&mut wasm, &mut enc).unwrap();
    enc.finish().unwrap();
    drop(enc);
    html.write_all(HTML_TAIL.as_bytes()).unwrap();
}

const HTML_HEAD: &'static str = r#"<html>
<head>
    <meta charset="UTF-8">
    <title>Rust WASM</title>
</head>
<body>
    <script type="module">
"#;

const HTML_MID: &'static str = r#"let wasm_base64 = ""#;

const HTML_TAIL: &'static str = r#"";
        let wasm_buffer = Uint8Array.from(atob(wasm_base64), c => c.charCodeAt(0)).buffer;
        window.onload = () => init(wasm_buffer);
    </script>
</body>
</html>"#;
