pub mod tiny_json {
    use std::fs;
    use std::fs::File;
    use std::io::ErrorKind;

    #[derive(Debug)]
    pub struct DB {
        path: String,
        file: File
    }

    pub fn open(path: &String) -> Result<DB, String> {
        match File::open(path) {
            Ok(handle) => Ok(DB { path: path.clone(), file: handle}),
            Err(e) => match e.kind() {
                ErrorKind::NotFound =>
                    match fs::create_dir(path) {
                        Err(_) => Err(String::from("failed to create data directory")),
                        Ok(_) => match File::open(path) {
                            Ok(file) => Ok(DB {file, path: path.clone()}),
                            Err(_) => Err(String::from("hello"))
                        }
                    },
                _ => Err(String::from("an error occurred attempting to open the data directory")),
            }
        }
    }

    impl DB {
        pub fn insert() {

        }
    }
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::fs::File;
    use super::*;

    #[test]
    fn test_open_created_directory_if_not_exists() {
        let test_data_path = String::from("./test-data");
        if does_file_exist(&test_data_path) {
            assert_eq!(true, fs::remove_dir_all(&test_data_path).is_ok())
        }
        assert_eq!(true, tiny_json::open(&test_data_path).is_ok());
        assert_eq!(true, does_file_exist(&test_data_path));
    }

    #[test]
    fn test_existing_directories_are_untouched() {
        let test_dir = String::from("./leave-me-alone");
        if !does_file_exist(&test_dir) {
            fs::create_dir_all("./leave-me-alone/file.json").expect("failed to create test dir");
        }

        assert_eq!(true, does_file_exist(&test_dir));
        assert_eq!(true, does_file_exist(&String::from("./leave-me-alone/file.json")));
        tiny_json::open(&test_dir).expect("failed to open database");
        assert_eq!(true, does_file_exist(&test_dir));
        assert_eq!(true, does_file_exist(&String::from("./leave-me-alone/file.json")));

        fs::remove_dir_all(&test_dir).expect("failed to cleanup test directory");
    }

    fn does_file_exist(filepath: &String) -> bool {
        File::open(filepath).is_ok()
    }
}