# Listmonk Terraform Provider

This Terraform provider allows you to manage templates in Listmonk.

## Installation

To use this provider, you need to have Terraform installed. You can download it from the official website: [Terraform Downloads](https://www.terraform.io/downloads.html).

## Usage

1. Create a new Terraform configuration file (e.g., `main.tf`) and add the necessary provider configuration:

    ```hcl
    provider "listmonk" {
      # Add your provider configuration here
    }
    ```

2. Define your template resources using the `listmonk_template` resource type:

    ```hcl
    resource "listmonk_template" "example_template" {
      # Add your template configuration here
    }
    ```

3. Run `terraform init` to initialize the provider and download the necessary dependencies.

4. Run `terraform apply` to create or update the resources defined in your configuration.

For more details on the resource configuration, refer to the documentation in the `/docs` folder.

## Limitations

Please note that this provider has limited capabilities and can only manage templates in Listmonk. For more advanced functionality, consider using the Listmonk API directly.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.
