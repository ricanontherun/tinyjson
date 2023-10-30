pub mod tiny_json {
    use std::fs;
    use std::fs::File;
    use std::io::ErrorKind;

    #[derive(Debug)]
    pub struct Handle {
        path: String,
        file: File
    }

    pub fn open(path: String) -> Result<Handle, String> {
        match File::open(path.clone()) {
            Ok(handle) => Ok(Handle{ path: path.clone(), file: handle}),
            Err(e) => match e.kind() {
                ErrorKind::NotFound =>
                    match fs::create_dir(path.clone()) {
                        Err(_) => Err(String::from("failed to create data directory")),
                        Ok(_) => match File::open(path.clone()) {
                            Ok(file) => Ok(Handle{file, path: path.clone()}),
                            Err(_) => Err(String::from("hello"))
                        }
                    },
                _ => Err(String::from("an error occurred attempting to open the data directory")),
            }
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
        let test_data_path = "./test-data";
        if does_file_exist(test_data_path) {
            assert_eq!(true, fs::remove_dir_all(test_data_path).is_ok())
        }
        assert_eq!(true, tiny_json::open(String::from(test_data_path)).is_ok());
        assert_eq!(true, does_file_exist(test_data_path));
    }

    #[test]
    fn test_existing_directories_are_untouched() {

    }

    fn does_file_exist(filepath: &str) -> bool {
        File::open(filepath).is_ok()
    }
}