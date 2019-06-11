package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/immesys/wave/consts"
	"github.com/immesys/wave/eapi"
	"github.com/immesys/wave/eapi/pb"
	"google.golang.org/grpc"
)

func resolveEntityNameOrHashOrFile(conn pb.WAVEClient, perspective *pb.Perspective, in string, msg string) (hash []byte) {
	f, err := ioutil.ReadFile(in)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error opening file %q: %v\n", in, err)
			os.Exit(1)
		}
		//Resolve as name/hash
		if len(in) == 48 && strings.Index(in, ".") == -1 {
			//Resolve as hash
			rv, err := base64.URLEncoding.DecodeString(in)
			if err != nil {
				fmt.Printf("bad base64: %q\n", in)
				os.Exit(1)
			}
			return rv
		}
		//Resolve as name
		if in == "wave" {
			//Hardcoded builtin PSET
			rv, _ := base64.URLEncoding.DecodeString(consts.WaveBuiltinPSET)
			return rv
		} else if in == "wavemq" {
			return []byte("\x1b\x20\x14\x33\x74\xb3\x2f\xd2\x74\x39\x54\xfe\x47\x86\xf6\xcf\x86\xd4\x03\x72\x0f\x5e\xc4\x42\x36\xb6\x58\xc2\x6a\x1e\x68\x0f\x6e\x01")
		}

		resp, err := conn.ResolveName(context.Background(), &pb.ResolveNameParams{
			Perspective: perspective,
			Name:        in,
		})
		if err != nil {
			fmt.Printf("could not resolve name: %v (%s)\n", err, msg)
			os.Exit(1)
		}
		if resp.Error != nil {
			fmt.Printf("could not resolve name %q: %s\n", in, resp.Error.Message)
			os.Exit(1)
		}
		return resp.Entity.Hash
	}
	//Resolve as file
	resp, err := conn.Inspect(context.Background(), &pb.InspectParams{
		Content: f,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	if resp.Error != nil {
		fmt.Printf("could not inspect file: %s\n", resp.Error.Message)
		os.Exit(1)
	}
	if resp.Entity != nil {
		return resp.Entity.Hash
	}
	fmt.Printf(msg)
	os.Exit(1)
	return nil
}

func getConn(agent string) pb.WAVEClient {
	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to agent: %v\n", err)
	}
	client := pb.NewWAVEClient(conn)
	return client
}

func getPerspective(file string, passphrase string, msg string) *pb.Perspective {
	if file != "" {
		pass := []byte(passphrase)
		if len(pass) == 0 {
			fmt.Printf("passphrase for entity secret: ")
			var err error
			pass, err = gopass.GetPasswdMasked()
			if err != nil {
				fmt.Printf("could not read passphrase: %v\n", err)
				os.Exit(1)
			}
		}
		pder := loadEntitySecretDER(file)
		perspective := &pb.Perspective{
			EntitySecret: &pb.EntitySecret{
				DER:        pder,
				Passphrase: pass,
			},
		}
		return perspective
	} else {
		fmt.Printf(msg)
		os.Exit(1)
		return nil
	}
}

func loadEntitySecretDER(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("could not read file %q: %v\n", filename, err)
		os.Exit(1)
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		fmt.Printf("file %q is not a PEM file\n", filename)
		os.Exit(1)
	}
	if block.Type != eapi.PEM_ENTITY_SECRET {
		fmt.Printf("PEM is not an entity secret\n")
		os.Exit(1)
	}
	return block.Bytes
}

func getNameFromHash(conn pb.WAVEClient, perspective *pb.Perspective, hash string) (name string) {
	bytehash, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		fmt.Printf("Hash %v was not base64 (%v)", hash, err)
		return hash
	}
	resp, err := conn.ResolveReverseName(context.Background(), &pb.ResolveReverseNameParams{
		Perspective: perspective,
		Hash:        bytehash,
	})
	if err != nil {
		fmt.Printf("Could not resolve name of hash %s %v\n", hash, err)
		return hash
	}
	if resp.Error != nil {
		fmt.Printf("Could not resolve name of hash %s %v\n", hash, resp.Error.Message)
		return hash
	}
	return resp.Name
}

func getHashFromName(conn pb.WAVEClient, perspective *pb.Perspective, name string) (hash string) {
	resp, err := conn.ResolveName(context.Background(), &pb.ResolveNameParams{
		Perspective: perspective,
		Name:        name,
	})
	if err != nil {
		fmt.Printf("Could not resolve hash for name %v\n", err)
		return name
	}
	if resp.Error != nil {
		fmt.Printf("Could not resolve hash for name %v\n", resp.Error.Message)
		return name
	}
	return base64.URLEncoding.EncodeToString(resp.Entity.Hash)
}
