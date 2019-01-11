use crate::model::output;

pub struct HTTPOutput {

}

impl HTTPOutput {
    pub fn new() -> HTTPOutput {
        return HTTPOutput {

        };
    }
}

impl output::Output for HTTPOutput {
    fn execute(&self) -> bool {
        return true;
    }
}
