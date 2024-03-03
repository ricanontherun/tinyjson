mod log {
    use std::fs::{File, OpenOptions};
    use std::io::{Seek, SeekFrom, Write};
    use std::net::Shutdown::Read;
    use std::os::unix::prelude::FileExt;
    use futures::AsyncWriteExt;
    use futures::future::err;
    use rev_lines::{RawRevLines, RevLines};

    // Log interface
    //  logs are append only, delimited files
    //  logs are read by obtaining an iterator
    //   records are read in reverse order.
    pub trait Log {
        fn write(&mut self, data: Vec<u8>) -> BooleanResult;

        fn first(&mut self) -> ReadResult;
    }

    pub struct LogIterator {
        inner_iter: RevLines<File>
    }

    impl LogIterator {
        pub fn from(file: File) -> LogIterator {
            let rev_lines = RevLines::new(file);
            LogIterator{inner_iter: rev_lines}
        }
    }

    impl Iterator for LogIterator {
        type Item = String;

        fn next(&mut self) -> Option<Self::Item> {
            match self.inner_iter.next() {
                None => None,
                Some(result) => match result {
                    Ok(s) => Some(s),
                    Err(re) => {
                        println!("failed to read next line, {}", re.to_string());
                        None
                    }
                }
            }
        }
    }

    pub fn open(file_path: &str) -> Result<Box<dyn Log>, &str> {
        let file_r = OpenOptions::new()
            .read(true)
            .write(true)
            .create(true)
            .truncate(false)
            .append(true)
            .open(file_path);

        if file_r.is_err() {
            return Err("failed to open/create log file");
        }

        let mut file = file_r.unwrap();

        return Ok(Box::new(DelimitedLogFile{file, record_delim: '\n'}));
    }

    struct DelimitedLogFile {
        file: std::fs::File,
        record_delim: char
    }

    impl Log for DelimitedLogFile {
        fn write(&mut self, data: Vec<u8>) -> BooleanResult {
            let mut record = data.clone();
            record.push(b'\n');
            match self.file.write(record.as_slice()) {
                Ok(_) => Ok(true),
                Err(e) => Err(e.to_string())
            }
        }

        fn first(&mut self ) -> ReadResult {
            let mut chars =vec![0;10];

            match self.file.read_exact_at(chars.as_mut_slice(), 0) {
                Ok(_) => ReadResult::Ok(String::from_utf8(chars).unwrap()),
                Err(e) => ReadResult::Err(e.to_string())
            }
        }
    }

    type BooleanResult = Result<bool, String>;
    type ReadResult = Result<String, String>;
}

#[cfg(test)]
mod tests {
    use super::*;

    // Utility methods.
    fn file_exists(path: &str) -> bool {
        match std::fs::File::open(path) {
            Ok(_) => true,
            _ => false
        }
    }

    fn file_does_not_exist(path: &str) -> bool {
        !file_exists(path)
    }

    fn clean_files() {
        let _ = std::fs::remove_file("test-file");
        assert!(file_does_not_exist("test-file"));
    }

    #[test]
    fn initializes_new_files() {
        clean_files();
        let log = log::open("test-file").unwrap();
        assert_file_exists("test-file");
    }

    fn assert_file_exists(file_path: &str) {
        assert!(std::fs::File::open(file_path).is_ok())
    }

    #[test]
    fn writes_to_disk() {
        clean_files();
        let mut log = log::open("test-file").expect("failed to open test-file");

        let write_result = log.write(Vec::from("some bytes"));
        assert!(write_result.is_ok());

        let write_result = log.write(Vec::from("some more bytes"));
        assert!(write_result.is_ok());

        let read_result = std::fs::read("test-file");
        assert!(read_result.is_ok());

        let file_bytes = read_result.unwrap();
        let file_contents = String::from_utf8(file_bytes);
        assert!(file_contents.is_ok());
        assert_eq!("some bytes\nsome more bytes\n", file_contents.unwrap())
    }

    #[test]
    fn read_records_in_reverse() {
        clean_files();

        let log = log::open("test-file").expect("failed to open test-file");


    }
}