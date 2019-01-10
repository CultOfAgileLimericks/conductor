pub trait Input {
    fn add_callback(&self, f: Box<Fn(Box<Input>) -> bool>);
}
