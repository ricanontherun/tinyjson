use std::io::Read;

struct DataFile {

}

struct Record {
    record: String
}

struct DataFileIterator {
    file: std::fs::File,
    // Reecord queue.
    // Its very possible, given the size of the byte buffer,
    // that multiple records can be read with a single call to next().
    // These records are queued up, and emitted individually.
    queue: Vec<Record>,
    buffer: Vec<u8>,
}

impl Iterator for DataFileIterator {
    type Item = Record;

    fn next(&mut self) -> Option<Self::Item> {

        None
    }
}

struct DataIterator {
    current: u8
}

impl Iterator for DataIterator {
    type Item = u8;

    fn next(&mut self) -> Option<Self::Item> {
        let next = self.current;
        self.current += 1;
        Some(next)
    }
}


#[cfg(test)]
mod tests {
    use crate::data::DataIterator;

    #[test]
    fn test() {
        let mut data_iterator = DataIterator{current: 0};
        assert_eq!(0, data_iterator.next().expect("expected data"));
        assert_eq!(1, data_iterator.next().expect("expected data"));
        assert_eq!(2, data_iterator.next().expect("expected data"));
        assert_eq!(3, data_iterator.next().expect("expected data"));
        assert_eq!(4, data_iterator.next().expect("expected data"));
    }
}
