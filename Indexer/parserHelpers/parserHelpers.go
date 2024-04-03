package parserhelpers

import (
	"fmt"
	"mail-indexer/constants"
	"mail-indexer/models"
	fetchzinc "mail-indexer/zincHelpers"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// adds the emails to allEmails array
func appendEmail(email *models.Email) {
	mu.Lock()
	defer mu.Unlock()
	allEmails = append(allEmails, *email)

	if len(allEmails) == constants.EMAILS_AMOUNT {
		fetchzinc.FetchZing(allEmails)
	}
}

// return s without extra white spaces
func RemoveExtraSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

func parseDateLines(line string, email *models.Email) {
	newLine := RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Date:")))
	firstParsedLine, firstErr := time.Parse(constants.REFERENCE_FIRST_DATE, newLine)
	if firstErr != nil {
		secondParsedLine, secondErr := time.Parse(constants.REFERENCE_SECOND_DATE, newLine)
		if secondErr != nil {
			email.DateSubEmail = newLine
			return
		}
		email.Date = secondParsedLine
		return
	}
	email.Date = firstParsedLine
}

// given a line, it adds the information to each email property
func ParseLineMessage(line string, email *models.Email, hasSubEmails *bool) {
	if strings.Contains(line, "=======================================") {
		email.Body += ""
	} else if strings.Contains(line, "Message-ID:") {
		email.MessageID = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:")))
	} else if strings.Contains(line, "Sent:") && !*hasSubEmails {
		email.Sent = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Sent:")))
	} else if strings.Contains(line, "Cc:") && !*hasSubEmails {
		email.Cc = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Cc:")))
	} else if strings.Contains(line, "X-From:") {
		email.XFrom = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-From:")))
	} else if strings.Contains(line, "X-To:") {
		email.XTo = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-To:")))
	} else if strings.Contains(line, "X-cc:") {
		email.XCc = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-cc:")))
	} else if strings.Contains(line, "X-bcc:") {
		email.XBcc = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-bcc:")))
	} else if strings.Contains(line, "X-Folder:") {
		email.XFolder = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-Folder:")))
	} else if strings.Contains(line, "X-Origin:") {
		email.XOrigin = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-Origin:")))
	} else if strings.Contains(line, "X-FileName:") {
		email.XFileName = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "X-FileName:")))
	} else if strings.Contains(line, "Date:") && !*hasSubEmails {
		parseDateLines(line, email)
	} else if strings.Contains(line, "From:") && !*hasSubEmails {
		email.From = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "From:")))
	} else if strings.Contains(line, "To:") && !*hasSubEmails {
		email.To += RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "To:"))) + " "
	} else if strings.Contains(line, "Subject:") && !*hasSubEmails {
		email.Subject = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Subject:")))
	} else if strings.Contains(line, "Mime-Version:") {
		email.MimeVersion = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:")))
	} else if strings.Contains(line, "Content-Type:") {
		email.ContentType = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:")))
	} else if strings.Contains(line, "Content-Transfer-Encoding:") {
		email.ContentTransferEncoding = RemoveExtraSpaces(strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:")))
	} else if strings.Contains(line, "-----Original Message-----") {
		*hasSubEmails = true
		email.Body += line + "\n"
	} else {
		email.Body += line + "\n"
	}
}

// returns a new id with prefix .JavaSubEmail
func generateSubEmailId(subEmailIndex int, emailId string) string {
	newId := fmt.Sprintf(".JavaSubEmail.%s>", strconv.Itoa(subEmailIndex))
	return strings.Replace(emailId, ">", newId, 1)
}

// checks if email contains subEmails and create new emails
func createSubEmails(email *models.Email) {
	if !strings.Contains(email.Body, "-----Original Message-----") {
		return
	}
	subEmails := strings.Split(email.Body, "-----Original Message-----")
	if len(subEmails) < 1 {
		return
	}
	for subEmailIndex := 1; subEmailIndex < len(subEmails); subEmailIndex++ {
		newEmail := models.Email{}
		hasSubEmails := false
		for _, line := range strings.Split(subEmails[subEmailIndex], "\n") {
			ParseLineMessage(line, &newEmail, &hasSubEmails)
		}
		newEmail.MessageID = generateSubEmailId(subEmailIndex, email.MessageID)
		appendEmail(&newEmail)
	}
	email.Body = subEmails[0]
}
