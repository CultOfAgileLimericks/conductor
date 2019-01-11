mod model;
mod plugin;

fn main() {
    let mut t = model::task::Task::new();
    let http_input = plugin::input::http::HTTPInput::new();
    let http_output = plugin::output::http::HTTPOutput::new();

    t.register_input(http_input);
    t.register_output(http_output);

    println!("{}", t);
}
