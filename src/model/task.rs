use std::fmt;

use super::input;
use super::output;

pub struct Task {
    pub inputs: Vec<Box<dyn input::Input>>,
    pub outputs: Vec<Box<dyn output::Output>>,
}

impl Task {
    pub fn new() -> Task {
        return Task { 
            inputs: Vec::new(),
            outputs: Vec::new(),
        }
    }

    pub fn register_input<T: input::Input + 'static>(&mut self, input_to_register: T) {
        input_to_register.add_callback(Box::new(|i: Box<input::Input>| {
            return true;
        }));
        self.inputs.push(Box::new(input_to_register));
    }

    pub fn register_output<T: output::Output + 'static>(&mut self, output_to_register: T) {
        self.outputs.push(Box::new(output_to_register));
    }
}

impl fmt::Display for Task {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "This is a task")
    }
}
