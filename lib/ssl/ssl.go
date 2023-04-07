package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/gobuffalo/packr/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Client() grpc.DialOption {
	box := packr.New("box1", "./client")

	// Get the string representation of a file, or an error if it doesn't exist:
	clientPem, err := box.Find("client.pem")
	if err != nil {
		logrus.Fatal(err)
	}
	clientKey, err := box.Find("client.key")
	if err != nil {
		logrus.Fatal(err)
	}
	cert, err := tls.X509KeyPair(clientPem, clientKey)
	if err != nil {
		logrus.Fatal(err)
	}
	//cert, err := tls.LoadX509KeyPair("./client/client.pem", "./client/client.key")
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	certPool := x509.NewCertPool()
	box2 := packr.New("box2", "./")
	caCrt, err := box2.Find("ca.crt")
	if err != nil {
		logrus.Fatal(err)
	}
	//ca, err := ioutil.ReadFile("./ca.crt")
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	// 尝试解析传入的pem证书，解析成功则会将其加入到certPool中，便于后面使用 - 貌似就是需要一个认证链的
	certPool.AppendCertsFromPEM(caCrt)

	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 校验客户端证书
		ServerName: "kmsm.qualstor.com",
		// 设置根证书集合，校验方式为ClientAuth中指定的模式
		RootCAs: certPool,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	return grpc.WithTransportCredentials(creds)
}

func Server() grpc.ServerOption {
	box := packr.New("box3", "./server")

	// Get the string representation of a file, or an error if it doesn't exist:
	clientPem, err := box.Find("server.pem")
	if err != nil {
		logrus.Fatal(err)
	}
	clientKey, err := box.Find("server.key")
	if err != nil {
		logrus.Fatal(err)
	}
	cert, err := tls.X509KeyPair(clientPem, clientKey)
	if err != nil {
		logrus.Fatal(err)
	}

	//cert, err := tls.LoadX509KeyPair("./server/server.pem", "./server/server.key")
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	certPool := x509.NewCertPool()
	box3 := packr.New("box4", "./")
	caCrt, err := box3.Find("ca.crt")
	if err != nil {
		logrus.Fatal(err)
	}
	//ca, err := ioutil.ReadFile("./ca.crt")
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	// 尝试解析传入的pem证书，解析成功则会将其加入到certPool中，便于后面使用 - 貌似就是需要一个认证链的
	certPool.AppendCertsFromPEM(caCrt)

	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须验证客户端证书，
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置验证客户端的根证书集合，校验方式为ClientAuth中指定的模式
		ClientCAs: certPool,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return grpc.Creds(creds)
}
