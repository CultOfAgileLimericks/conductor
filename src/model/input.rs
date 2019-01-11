pub trait Input {
    fn add_callback(&mut self, f: Box<Fn(Box<Input>) -> bool>);
}
