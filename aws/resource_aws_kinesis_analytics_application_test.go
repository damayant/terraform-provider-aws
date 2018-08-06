package aws

import (
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalytics"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSKinesisAnalyticsApplication_basic(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKinesisAnalyticsApplication_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "code", "testCode\n"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_update(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKinesisAnalyticsApplication_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
				),
			},
			{
				Config: testAccKinesisAnalyticsApplication_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "code", "testCode2\n"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_addCloudwatchLoggingOptions(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_basic(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_cloudwatchLoggingOptions(rInt, "testStream")
	streamRe := regexp.MustCompile(fmt.Sprintf("^arn:.*:log-stream:testAcc-testStream$"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "1"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "cloudwatch_logging_options.#", "1"),
					resource.TestMatchResourceAttr(resName, "cloudwatch_logging_options.0.log_stream", streamRe),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_updateCloudwatchLoggingOptions(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_cloudwatchLoggingOptions(rInt, "testStream")
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_cloudwatchLoggingOptions(rInt, "testStream2")
	beforeRe := regexp.MustCompile(fmt.Sprintf("^arn:.*:log-stream:testAcc-testStream$"))
	afterRe := regexp.MustCompile(fmt.Sprintf("^arn:.*:log-stream:testAcc-testStream2$"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "cloudwatch_logging_options.#", "1"),
					resource.TestMatchResourceAttr(resName, "cloudwatch_logging_options.0.log_stream", beforeRe),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "cloudwatch_logging_options.#", "1"),
					resource.TestMatchResourceAttr(resName, "cloudwatch_logging_options.0.log_stream", afterRe),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_inputsKinesisStream(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_inputsKinesisStream(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.name_prefix", "test_prefix"),
					resource.TestCheckResourceAttr(resName, "inputs.0.kinesis_stream.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.parallelism.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_columns.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.0.mapping_parameters.0.json.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_inputsAdd(t *testing.T) {
	var before, after kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_basic(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_inputsKinesisStream(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &before),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.#", "0"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &after),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "inputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.name_prefix", "test_prefix"),
					resource.TestCheckResourceAttr(resName, "inputs.0.kinesis_stream.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.parallelism.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_columns.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.0.mapping_parameters.0.json.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_inputsUpdateKinesisStream(t *testing.T) {
	var before, after kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_inputsKinesisStream(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_inputsUpdateKinesisStream(rInt, "testStream")
	streamRe := regexp.MustCompile(fmt.Sprintf("^arn:.*:stream/testAcc-testStream$"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &before),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.name_prefix", "test_prefix"),
					resource.TestCheckResourceAttr(resName, "inputs.0.parallelism.0.count", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.0.mapping_parameters.0.json.#", "1"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &after),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "inputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "inputs.0.name_prefix", "test_prefix2"),
					resource.TestMatchResourceAttr(resName, "inputs.0.kinesis_stream.0.resource", streamRe),
					resource.TestCheckResourceAttr(resName, "inputs.0.parallelism.0.count", "2"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_columns.0.name", "test2"),
					resource.TestCheckResourceAttr(resName, "inputs.0.schema.0.record_format.0.mapping_parameters.0.csv.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_outputsKinesisStream(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_outputsKinesisStream(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.name", "test_name"),
					resource.TestCheckResourceAttr(resName, "outputs.0.kinesis_stream.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.0.record_format_type", "JSON"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_outputsAdd(t *testing.T) {
	var before, after kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_basic(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_outputsKinesisStream(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &before),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.#", "0"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &after),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "outputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.name", "test_name"),
					resource.TestCheckResourceAttr(resName, "outputs.0.kinesis_stream.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_outputsUpdateKinesisStream(t *testing.T) {
	var before, after kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_outputsKinesisStream(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_outputsUpdateKinesisStream(rInt, "testStream")
	streamRe := regexp.MustCompile(fmt.Sprintf("^arn:.*:stream/testAcc-testStream$"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &before),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.name", "test_name"),
					resource.TestCheckResourceAttr(resName, "outputs.0.kinesis_stream.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.0.record_format_type", "JSON"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &after),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "outputs.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.name", "test_name2"),
					resource.TestCheckResourceAttr(resName, "outputs.0.kinesis_stream.#", "1"),
					resource.TestMatchResourceAttr(resName, "outputs.0.kinesis_stream.0.resource", streamRe),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "outputs.0.schema.0.record_format_type", "CSV"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_referenceDataSource(t *testing.T) {
	var application kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_referenceDataSource(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &application),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.#", "1"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.0.schema.#", "1"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.0.schema.0.record_columns.#", "1"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.0.schema.0.record_format.#", "1"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.0.schema.0.record_format.0.mapping_parameters.0.json.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSKinesisAnalyticsApplication_referenceDataSourceUpdate(t *testing.T) {
	var before, after kinesisanalytics.ApplicationDetail
	resName := "aws_kinesis_analytics_application.test"
	rInt := acctest.RandInt()
	firstStep := testAccKinesisAnalyticsApplication_prereq(rInt)
	secondStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_referenceDataSource(rInt)
	thirdStep := testAccKinesisAnalyticsApplication_prereq(rInt) + testAccKinesisAnalyticsApplication_referenceDataSourceUpdate(rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisAnalyticsApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstStep,
				Check: resource.ComposeTestCheckFunc(
					fulfillSleep(),
				),
			},
			{
				Config: secondStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &before),
					resource.TestCheckResourceAttr(resName, "version", "2"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.#", "1"),
				),
			},
			{
				Config: thirdStep,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisAnalyticsApplicationExists(resName, &after),
					resource.TestCheckResourceAttr(resName, "version", "3"),
					resource.TestCheckResourceAttr(resName, "reference_data_sources.#", "1"),
				),
			},
		},
	})
}

func testAccCheckKinesisAnalyticsApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_kinesis_analytics_application" {
			continue
		}
		conn := testAccProvider.Meta().(*AWSClient).kinesisanalyticsconn
		describeOpts := &kinesisanalytics.DescribeApplicationInput{
			ApplicationName: aws.String(rs.Primary.Attributes["name"]),
		}
		resp, err := conn.DescribeApplication(describeOpts)
		if err == nil {
			if resp.ApplicationDetail != nil && *resp.ApplicationDetail.ApplicationStatus != kinesisanalytics.ApplicationStatusDeleting {
				return fmt.Errorf("Error: Application still exists")
			}
		}
		return nil
	}
	return nil
}

func testAccCheckKinesisAnalyticsApplicationExists(n string, application *kinesisanalytics.ApplicationDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kinesis Analytics Application ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).kinesisanalyticsconn
		describeOpts := &kinesisanalytics.DescribeApplicationInput{
			ApplicationName: aws.String(rs.Primary.Attributes["name"]),
		}
		resp, err := conn.DescribeApplication(describeOpts)
		if err != nil {
			return err
		}

		*application = *resp.ApplicationDetail

		return nil
	}
}

func testAccKinesisAnalyticsApplication_basic(rInt int) string {
	return fmt.Sprintf(`
resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"
}
`, rInt)
}

func testAccKinesisAnalyticsApplication_update(rInt int) string {
	return fmt.Sprintf(`
resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode2\n"
}
`, rInt)
}

func testAccKinesisAnalyticsApplication_cloudwatchLoggingOptions(rInt int, streamName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = "testAcc-%d"
}

resource "aws_cloudwatch_log_stream" "test" {
  name = "testAcc-%s"
  log_group_name = "${aws_cloudwatch_log_group.test.name}"
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"

  cloudwatch_logging_options {
    log_stream = "${aws_cloudwatch_log_stream.test.arn}"
    role = "${aws_iam_role.test.arn}"
  }
}
`, rInt, streamName, rInt)
}

func testAccKinesisAnalyticsApplication_inputsKinesisStream(rInt int) string {
	return fmt.Sprintf(`
resource "aws_kinesis_stream" "test" {
  name = "testAcc-%d"
  shard_count = 1
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"

  inputs {
    name_prefix = "test_prefix"
    kinesis_stream {
      resource = "${aws_kinesis_stream.test.arn}"
      role = "${aws_iam_role.test.arn}"
    }
    parallelism {
      count = 1
    }
    schema {
      record_columns {
        mapping = "$.test"
        name = "test"
        sql_type = "VARCHAR(8)"
      }
      record_encoding = "UTF-8"
      record_format {
        mapping_parameters {
          json {
            record_row_path = "$"
          }
        }
      }
    }
  }
}
`, rInt, rInt)
}

func testAccKinesisAnalyticsApplication_inputsUpdateKinesisStream(rInt int, streamName string) string {
	return fmt.Sprintf(`
resource "aws_kinesis_stream" "test" {
  name = "testAcc-%s"
  shard_count = 1
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"

  inputs {
    name_prefix = "test_prefix2"
    kinesis_stream {
      resource = "${aws_kinesis_stream.test.arn}"
      role = "${aws_iam_role.test.arn}"
    }
    parallelism {
      count = 2
    }
    schema {
      record_columns {
        mapping = "$.test2"
        name = "test2"
        sql_type = "VARCHAR(8)"
      }
      record_encoding = "UTF-8"
      record_format {
        mapping_parameters {
          csv {
            record_column_delimiter = ","
            record_row_delimiter = "\n"
          }
        }
      }
    }
  }
}
`, streamName, rInt)
}

func testAccKinesisAnalyticsApplication_outputsKinesisStream(rInt int) string {
	return fmt.Sprintf(`
resource "aws_kinesis_stream" "test" {
  name = "testAcc-%d"
  shard_count = 1
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"

  outputs {
    name = "test_name"
    kinesis_stream {
      resource = "${aws_kinesis_stream.test.arn}"
      role = "${aws_iam_role.test.arn}"
    }
    schema {
      record_format_type = "JSON"
    }
  }
}
`, rInt, rInt)
}

func testAccKinesisAnalyticsApplication_outputsUpdateKinesisStream(rInt int, streamName string) string {
	return fmt.Sprintf(`
resource "aws_kinesis_stream" "test" {
  name = "testAcc-%s"
  shard_count = 1
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"
  code = "testCode\n"

  outputs {
    name = "test_name2"
    kinesis_stream {
      resource = "${aws_kinesis_stream.test.arn}"
      role = "${aws_iam_role.test.arn}"
    }
    schema {
      record_format_type = "CSV"
    }
  }
}
`, streamName, rInt)
}

func testAccKinesisAnalyticsApplication_referenceDataSource(rInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = "testacc-%d"
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"

  reference_data_sources {
    table_name = "test_table"
    s3 {
      bucket = "${aws_s3_bucket.test.arn}"
      file_key = "test_file_key"
      role = "${aws_iam_role.test.arn}"
    }
    schema {
      record_columns {
        mapping = "$.test"
        name = "test"
        sql_type = "VARCHAR(8)"
      }
      record_encoding = "UTF-8"
      record_format {
        mapping_parameters {
          json {
            record_row_path = "$"
          }
        }
      }
    }
  }
}
`, rInt, rInt)
}

func testAccKinesisAnalyticsApplication_referenceDataSourceUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = "testacc2-%d"
}

resource "aws_kinesis_analytics_application" "test" {
  name = "testAcc-%d"

  reference_data_sources {
    table_name = "test_table2"
    s3 {
      bucket = "${aws_s3_bucket.test.arn}"
      file_key = "test_file_key"
      role = "${aws_iam_role.test.arn}"
    }
    schema {
      record_columns {
        mapping = "$.test2"
        name = "test2"
        sql_type = "VARCHAR(8)"
      }
      record_encoding = "UTF-8"
      record_format {
        mapping_parameters {
          csv {
            record_column_delimiter = ","
            record_row_delimiter = "\n"
          }
        }
      }
    }
  }
}
`, rInt, rInt)
}

func fulfillSleep() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Print("[DEBUG] Test: Sleep to allow IAM to propagate")
		time.Sleep(30 * time.Second)
		return nil
	}
}

// this is used to set up the IAM role
func testAccKinesisAnalyticsApplication_prereq(rInt int) string {
	return fmt.Sprintf(`
data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["kinesisanalytics.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "test" {
  name = "testAcc-%d"
  assume_role_policy = "${data.aws_iam_policy_document.test.json}" 
}
`, rInt)
}
