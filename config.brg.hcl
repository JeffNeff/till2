source "awssqs" "my_queue" {
  arn = "arn:aws:sqs:us-east-2:123456789012:triggermesh"
  to = router.my_router
}

function "my_function" {
  repo = "github.com/acme/my-function"
  reply_to = router.my_router
}

router "content_based" "my_router" {

  route {
    attributes = {
      type = "com.amazon.sqs.message"
    }
    to = transformer.my_transformation
  }

  route {
    attributes = {
      type = "corp.acme.my.processing"
    }
    to = target.my_kafka_topic
  }

}

transformer "bumblebee" "my_transformation" {

  context {
    operation "store" {
      path {
        key = "$id"
        value = "id"
      }
    }
    operation "add" {
      path {
        key = "id"
        value = "${person}-${id}"
      }
    }
  }

  data {
    operation "store" {
      path {
        key = "$person"
        value = "Alice"
      }
    }
    operation "add" {
      path {
        key = "event.ID"
        value = "$id"
      }
    }
  }

  to = function.my_function
}

target "kafka" "my_kafka_topic" {
  topic = "myapp"
}
