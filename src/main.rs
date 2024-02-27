use std::fs::File;
use std::io;
use std::io::Read;
use futures::stream::{self, StreamExt};

#[derive(Debug)]
struct Record {
    data: String,
    start: u128,
}

struct DataFile {
    path: String
}

struct DataFileReader {
    data_file: DataFile
}

async fn main() -> io::Result<()> {
    // Thanks borrow checker, this will be closed when its last borrow drops the "lease"
    let data_file = DataFile{path: String::from("./people")};

    // Data file storage
    //  primary data files (source of truth)
    //  index files (point to a particular record in a data file)

    // Full seek (lets say looking for a particular ID)
    // We must read each character, looking for newlines.
    let mut file = File::open("./people")
        .expect("unable to open file 'people'");

    let mut start = 0;
    let mut i = 0;
    let mut record = String::new();
    let mut records: Vec<Record> = Vec::new();
    // A buffer which stores 10 bytes.
    let mut buffer = [0; 10];
    while file.read(&mut buffer[..])? != 0 {
        // If we've found a record delimiter (newline), cut the record.
        for c in buffer {
            i += 1;
            if c == '\n' as u8 {
                // Cut the record
                records.push(Record{data: record.clone(), start });
                record.clear();
                start = i + 1;
            } else {
                record.push(c as char);
            }
        }
    }

    for record in records {
        println!("record: {:?}", record);
        println!("starts={}, ends={}", record.start, (record.start + record.data.len() as u128))
    }

    Ok(())
}

#[cfg(test)]
mod test {
    #[test]
    fn reads_records_from_file() {

    }
    // Empty file, no records
    // Single data, single record
    // multiple delimiters in between record, shouldn't result in a record.
}