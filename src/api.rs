pub mod tiny_json {
    use std::fs;

    #[derive(Debug)]
    pub struct Handle {
        data_path: String,
    }

    pub fn open(path: String) -> Result<Handle, String> {
        match fs::create_dir(path.clone()) {
            Err(_) => Err(String::from("failed to create file")),
            _ => {
                Ok(Handle { data_path: path.clone() })
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
        assert_eq!(true, does_file_exist(test_data_path))
    }

    fn does_file_exist(filepath: &str) -> bool {
        File::open(filepath).is_ok()
    }
}