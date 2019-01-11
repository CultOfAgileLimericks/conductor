use crate::model::input;

pub struct HTTPInput {
    callbacks: Vec<Box<Fn(Box<input::Input>) -> bool>>,
}

impl HTTPInput {
    pub fn new() -> HTTPInput {
        return HTTPInput {
            callbacks: Vec::new(),
        };
    }
}

impl input::Input for HTTPInput {
    fn add_callback(&mut self, f: Box<Fn(Box<input::Input>) -> bool>) {
        self.callbacks.push(f);
    }
}
