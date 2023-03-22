use gloo::timers::future::sleep;
use instant::Instant;
use rand::{rngs::SmallRng, Rng, SeedableRng};
use std::time::Duration;
use wasm_bindgen::{Clamped, JsCast};
use web_sys::{CanvasRenderingContext2d, HtmlCanvasElement, ImageData};

fn main() {
    std::panic::set_hook(Box::new(console_error_panic_hook::hook));
    wasm_logger::init(wasm_logger::Config::new(log::Level::Info));
    wasm_bindgen_futures::spawn_local(run())
}

async fn run() {
    // Magnitude of side length
    let mag: usize = 8;
    // Side length (e.g. 256 for mag=8)
    let len: usize = 1 << mag;
    let len_mask = len - 1;
    // Total number of sites
    let size: usize = len * len;
    let size_mask = size - 1;

    // Precalculate probabilities for delta in range [-4, ..., +4]
    let mut p = [0f64; 9];
    for i in 0..p.len() {
        let beta = f64::ln(1.0 + f64::sqrt(2.0)) / 2.0; // 0.44068679350977147
        p[i] = f64::exp(-2.0 * beta * (i as f64 - 4.0))
    }

    let ctx = create_canvas(len);

    let mut state = vec![1i8; size];

    let mut rng = SmallRng::seed_from_u64(0);
    let start = Instant::now();
    let mut last_draw = Instant::now();
    let mut sweeps = 0;
    loop {
        for off in 0..3 {
            for center in (off..size).step_by(3) {
                let col = center & len_mask;
                let row_off = center - col;
                let (down, up) = (center + len & size_mask, center + size - len & size_mask);
                let (right, left) = (row_off + (col + 1 & len_mask), row_off + (col + len - 1 & len_mask));

                let delta = state[center] * (state[left] + state[right] + state[up] + state[down]);
                if delta <= 0 || rng.gen::<f64>() < p[(delta + 4) as usize] {
                    state[center] *= -1
                }
            }
        }
        sweeps += 1;
        if last_draw.elapsed() > Duration::from_millis(100) {
            last_draw = Instant::now();
            log::info!("sweep rate: {}/s", sweeps as f64 / start.elapsed().as_secs_f64());
            draw(len, &state, &ctx);
            // Give the browser a chance to present
            sleep(Duration::from_millis(1)).await
        }
    }
}

fn create_canvas(len: usize) -> CanvasRenderingContext2d {
    let document = web_sys::window().unwrap().document().unwrap();
    let canvas: HtmlCanvasElement = document.create_element("canvas").unwrap().dyn_into().unwrap();
    canvas.style().set_property("image-rendering", "pixelated").unwrap();
    canvas.set_height(len as u32);
    canvas.set_width(len as u32);
    document.body().unwrap().append_child(&canvas).unwrap();
    return canvas.get_context("2d").unwrap().unwrap().dyn_into().unwrap();
}

fn draw(len: usize, state: &Vec<i8>, ctx: &CanvasRenderingContext2d) {
    let mut pixels = vec![0u8; len * len * 4];
    for i in 0..len * len {
        let pixel = &mut pixels[i * 4..i * 4 + 4];
        pixel[3] = 255; // alpha
        if state[i] < 0 {
            pixel[0] = 255; // red
        } else {
            pixel[2] = 255; // blue
        }
    }
    let img = ImageData::new_with_u8_clamped_array(Clamped(&pixels), len as u32).unwrap();
    ctx.put_image_data(&img, 0.0, 0.0).unwrap();
}
